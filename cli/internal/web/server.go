// Package web provides a small read-only HTTP dashboard that renders the
// pressctl configuration (servers, sites, domains) as a styled web page.
package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pressctl/cli/internal/config"
	"github.com/pressctl/cli/pkg/models"
)

//go:embed dashboard.html
var assets embed.FS

var tmpl = template.Must(template.ParseFS(assets, "dashboard.html"))

// viewData is the top-level template payload.
type viewData struct {
	ConfigPath  string
	Now         string
	ServerCount int
	SiteCount   int
	DomainCount int
	SSLCount    int
	Servers     []serverView
	Sites       []siteView // flat list of every site across all servers
}

type serverView struct {
	Name        string
	Hostname    string
	IP          string
	Status      string
	Stack       string
	PHPVersion  string
	Provisioned string
	SSHCommand  string
	SiteCount   int
}

type siteView struct {
	ServerName    string
	SiteID        string
	PrimaryDomain string
	AdminURL      string
	ExtraDomains  string
	SSLEnabled    bool
	SSLState      string // "ok" | "warn" | "err" | "none"
	SSLLabel      string // human label, e.g. "valid · 42d" or "expired"
	AdminUser     string
	PHPVersion    string
	Database      string
	BackupLabel   string
	Notes         string
	Created       string
}

// Serve loads the config, starts a local HTTP server on the given port, and
// optionally opens the default browser. It blocks until ctx is cancelled.
// A port of 0 selects a random free port.
func Serve(ctx context.Context, mgr *config.Manager, port int, open bool) error {
	configPath := mgr.GetConfigPath()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Reload on every request so the page reflects external edits.
		cfg, err := mgr.Load()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load config: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, buildView(configPath, cfg)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("failed to bind port: %w", err)
	}

	url := fmt.Sprintf("http://%s", ln.Addr().String())
	srv := &http.Server{Handler: mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("pressctl dashboard running at %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")
	if open {
		openBrowser(url)
	}

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func buildView(configPath string, cfg *config.Config) viewData {
	v := viewData{
		ConfigPath: configPath,
		Now:        time.Now().Format("Jan 2, 2006 15:04:05"),
	}

	for _, s := range cfg.Servers {
		sv := serverView{
			Name:        s.Name,
			Hostname:    s.Hostname,
			IP:          s.IP,
			Status:      s.Status,
			Stack:       s.EffectiveStack(),
			PHPVersion:  orDash(s.PHPVersion),
			Provisioned: formatDatePtr(s.ProvisionedAt),
			SSHCommand:  sshCommand(s),
			SiteCount:   len(s.Sites),
		}
		for _, site := range s.Sites {
			sview := buildSite(site)
			sview.ServerName = s.Name
			v.Sites = append(v.Sites, sview)
			v.SiteCount++
			v.DomainCount += len(site.Domains)
			for _, d := range site.Domains {
				if d.SSLEnabled {
					v.SSLCount++
				}
			}
		}
		v.Servers = append(v.Servers, sv)
	}
	v.ServerCount = len(cfg.Servers)
	return v
}

func buildSite(site models.Site) siteView {
	var extras []string
	var primary models.Domain
	for _, d := range site.Domains {
		if d.Domain == site.PrimaryDomain {
			primary = d
			continue
		}
		extras = append(extras, d.Domain)
	}

	sslState, sslLabel := sslStatus(primary)

	return siteView{
		SiteID:        site.SiteID,
		PrimaryDomain: site.PrimaryDomain,
		AdminURL:      "https://" + site.PrimaryDomain + "/wp-admin",
		ExtraDomains:  strings.Join(extras, ", "),
		SSLEnabled:    primary.SSLEnabled,
		SSLState:      sslState,
		SSLLabel:      sslLabel,
		AdminUser:     site.AdminUser,
		PHPVersion:    orDash(site.PHPVersion),
		Database:      orDash(site.Database.Name),
		BackupLabel:   backupLabel(site.Metadata),
		Notes:         site.Notes,
		Created:       site.CreatedAt.Format("Jan 2, 2006"),
	}
}

// sslStatus returns a coloring state and a human label for a domain's
// certificate, factoring in days-until-expiry when known.
func sslStatus(d models.Domain) (state, label string) {
	if !d.SSLEnabled {
		return "none", "off"
	}
	if d.SSLExpiresAt == nil {
		return "ok", "on"
	}
	days := int(time.Until(*d.SSLExpiresAt).Hours() / 24)
	switch {
	case days < 0:
		return "err", "expired"
	case days <= 14:
		return "warn", fmt.Sprintf("%dd left", days)
	default:
		return "ok", fmt.Sprintf("%dd left", days)
	}
}

func backupLabel(m models.Metadata) string {
	if !m.BackupEnabled {
		return "off"
	}
	if m.LastBackup != nil {
		return "on · " + m.LastBackup.Format("Jan 2")
	}
	return "on"
}

func sshCommand(s models.Server) string {
	user := s.SSH.User
	if user == "" {
		user = "pressctl"
	}
	port := s.SSH.Port
	if port == 0 {
		port = 22
	}
	return fmt.Sprintf("ssh -p %d %s@%s", port, user, s.IP)
}

func formatDatePtr(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format("Jan 2, 2006")
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

// openBrowser tries to open url in the user's default browser. Failures are
// non-fatal; the URL is already printed for manual access.
func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd, args = "cmd", []string{"/c", "start"}
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	_ = exec.Command(cmd, args...).Start()
}
