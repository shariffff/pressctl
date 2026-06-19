package web

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/pressctl/cli/internal/config"
	"github.com/pressctl/cli/pkg/models"
)

func TestBuildViewAndRender(t *testing.T) {
	created := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	expires := time.Now().Add(40 * 24 * time.Hour)
	cfg := &config.Config{
		Servers: []models.Server{{
			Name:          "web-1",
			Hostname:      "web1.example.com",
			IP:            "203.0.113.10",
			Status:        "provisioned",
			Stack:         "frankenphp",
			SSH:           models.SSHConfig{User: "pressctl", Port: 2222},
			ProvisionedAt: &created,
			Sites: []models.Site{{
				SiteID:        "blogsite",
				PrimaryDomain: "blog.example.com",
				AdminUser:     "admin",
				PHPVersion:    "8.3",
				CreatedAt:     created,
				Database:      models.Database{Name: "blogsite_db"},
				Metadata:      models.Metadata{BackupEnabled: true},
				Notes:         "migrate before July",
				Domains: []models.Domain{
					{Domain: "blog.example.com", SSLEnabled: true, SSLExpiresAt: &expires},
					{Domain: "www.blog.example.com", SSLEnabled: true},
				},
			}},
		}},
	}

	v := buildView("/tmp/pressctl.yaml", cfg)
	if v.ServerCount != 1 || v.SiteCount != 1 || v.DomainCount != 2 || v.SSLCount != 2 {
		t.Fatalf("unexpected counts: %+v", v)
	}
	srv := v.Servers[0]
	if srv.SSHCommand != "ssh -p 2222 pressctl@203.0.113.10" {
		t.Fatalf("ssh command = %q", srv.SSHCommand)
	}
	if srv.SiteCount != 1 {
		t.Fatalf("server site count = %d", srv.SiteCount)
	}
	site := v.Sites[0]
	if site.ServerName != "web-1" {
		t.Fatalf("site server name = %q", site.ServerName)
	}
	if site.ExtraDomains != "www.blog.example.com" {
		t.Fatalf("extra domains = %q", site.ExtraDomains)
	}
	if !site.SSLEnabled || site.SSLState != "ok" || !strings.Contains(site.SSLLabel, "d left") {
		t.Fatalf("ssl = %+v", site)
	}
	if site.Database != "blogsite_db" || site.BackupLabel != "on" || site.Notes != "migrate before July" {
		t.Fatalf("site detail fields = %+v", site)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, v); err != nil {
		t.Fatalf("template execute: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"web-1", "blog.example.com", "blogsite", "Copy SSH", "wp-admin", "blogsite_db", "migrate before July"} {
		if !strings.Contains(out, want) {
			t.Errorf("rendered output missing %q", want)
		}
	}
}
