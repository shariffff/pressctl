package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/pressctl/cli/internal/ansible"
	"github.com/pressctl/cli/pkg/models"
	"github.com/spf13/cobra"
)

// siteConnectCmd represents the site connect command
var siteConnectCmd = &cobra.Command{
	Use:     "connect",
	Aliases: []string{"connect-wordsail"},
	Short:   "Connect a site to WordSail",
	Long: `Connect a site's WordSail connector plugin by writing the connection key
to the site's wp-config.php as WORDSAIL_CONNECTION_KEY. Re-running with a new
key replaces the existing one.

The key can be given with --key or the WORDSAIL_CONNECTION_KEY environment
variable (preferred in scripts, since command-line arguments are visible to
other users via the process list).

Examples:
  # Interactive mode
  press site connect

  # Non-interactive mode (for automation/AI agents)
  press site connect --server myserver --site mysite --key KEY
  WORDSAIL_CONNECTION_KEY=KEY press site connect --server myserver --site mysite`,
	Run: func(cmd *cobra.Command, args []string) {
		_, cfg := ensureConfig()

		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			key = os.Getenv("WORDSAIL_CONNECTION_KEY")
		}
		key = strings.TrimSpace(key)

		if (serverName == "") != (siteName == "") {
			color.Red("Error: --server and --site must be provided together")
			os.Exit(1)
		}

		// Interactive site selection
		if serverName == "" {
			type siteOption struct {
				ServerName string
				Site       models.Site
			}

			var siteOptions []siteOption
			for _, server := range cfg.Servers {
				for _, site := range server.Sites {
					siteOptions = append(siteOptions, siteOption{ServerName: server.Name, Site: site})
				}
			}

			if len(siteOptions) == 0 {
				fmt.Println("No sites available.")
				fmt.Println("Create a site with: press site create")
				return
			}

			opts := make([]huh.Option[int], len(siteOptions))
			for i, o := range siteOptions {
				opts[i] = huh.NewOption(
					fmt.Sprintf("%s on %s (%s)", o.Site.PrimaryDomain, o.ServerName, o.Site.SiteID), i)
			}
			var selected int
			if err := huh.NewSelect[int]().
				Title("Select site to connect").
				Options(opts...).
				Value(&selected).
				Run(); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			serverName = siteOptions[selected].ServerName
			siteName = siteOptions[selected].Site.SiteID
		}

		// Prompt for the key (masked) if not supplied
		if key == "" {
			if err := huh.NewInput().
				Title("WordSail connection key").
				EchoMode(huh.EchoModePassword).
				Value(&key).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("connection key is required")
					}
					return nil
				}).
				Run(); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			key = strings.TrimSpace(key)
		}

		// Find the server and site
		var targetServer *models.Server
		var targetSite *models.Site
		for i := range cfg.Servers {
			if cfg.Servers[i].Name != serverName {
				continue
			}
			targetServer = &cfg.Servers[i]
			for j := range cfg.Servers[i].Sites {
				if cfg.Servers[i].Sites[j].SiteID == siteName {
					targetSite = &cfg.Servers[i].Sites[j]
					break
				}
			}
			break
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", serverName)
			os.Exit(1)
		}
		if targetSite == nil {
			color.Red("Error: Site '%s' not found on server '%s'", siteName, serverName)
			os.Exit(1)
		}

		extraVars := map[string]interface{}{
			"site_id":                 targetSite.SiteID,
			"site_domain":             targetSite.PrimaryDomain,
			"wordsail_connection_key": key,
		}

		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Connecting %s to WordSail", targetSite.PrimaryDomain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("playbooks/connect_wordsail.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Failed to connect site: %v", err)
			os.Exit(1)
		}

		fmt.Println()
		color.Green("✓ WORDSAIL_CONNECTION_KEY set for '%s'", targetSite.PrimaryDomain)
	},
}

func init() {
	siteCmd.AddCommand(siteConnectCmd)

	siteConnectCmd.Flags().String("server", "", "Server name")
	siteConnectCmd.Flags().String("site", "", "Site ID")
	siteConnectCmd.Flags().String("key", "", "WordSail connection key (or set WORDSAIL_CONNECTION_KEY)")
}
