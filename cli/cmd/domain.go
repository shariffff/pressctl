package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/pressctl/cli/internal/ansible"
	"github.com/pressctl/cli/internal/prompt"
	"github.com/pressctl/cli/internal/state"
	"github.com/pressctl/cli/internal/utils"
	"github.com/pressctl/cli/pkg/models"
	"github.com/spf13/cobra"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage domains and SSL certificates",
	Long:  `Add domains to sites, remove domains, and issue SSL certificates.`,
}

// domainAddCmd represents the domain add command
var domainAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Short:   "Add a domain to a site",
	Long: `Add a new domain to an existing WordPress site and optionally issue an SSL certificate.

Examples:
  # Interactive mode
  press domain add

  # Non-interactive mode (for automation/AI agents)
  press domain add --server myserver --site mysite --domain www.example.com --ssl`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		var input *prompt.DomainAddInput

		// Check for non-interactive mode
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")

		if serverName != "" && siteName != "" && domain != "" {
			// Non-interactive mode
			issueSSL, _ := cmd.Flags().GetBool("ssl")
			input = &prompt.DomainAddInput{
				ServerName: serverName,
				SiteID:     siteName,
				Domain:     domain,
				IssueSSL:   issueSSL,
			}
		} else if serverName != "" || siteName != "" || domain != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are all required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - get input from prompts
			var err error
			input, err = prompt.PromptDomainAdd(cfg.Servers)
			if err != nil {
				outputError(cmd, "Failed to get domain details", err)
				os.Exit(1)
			}
		}

		// The domain becomes root file paths and nginx config on the server.
		input.Domain = strings.ToLower(strings.TrimSpace(input.Domain))
		if err := utils.ValidateDomain(input.Domain); err != nil {
			outputError(cmd, "Invalid domain", fmt.Errorf("%q: %w", input.Domain, err))
			os.Exit(1)
		}

		// A hostname can belong to one site only: a second vhost or certificate
		// for it would take traffic or the certificate from the other site.
		for i := range cfg.Servers {
			if owner := utils.FindSiteByDomain(&cfg.Servers[i], input.Domain); owner != nil {
				outputError(cmd, "Domain in use", fmt.Errorf("%s is already used by site '%s' on server '%s'", input.Domain, owner.SiteID, cfg.Servers[i].Name))
				os.Exit(1)
			}
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		// Prepare extra vars for Ansible
		alias, _ := cmd.Flags().GetBool("alias")
		extraVars := map[string]interface{}{
			"operation": "add_domain",
			"domain":    input.Domain,
			"site_id":   input.SiteID,
		}
		if alias {
			vars, err := aliasVars(targetServer, input.SiteID, input.Domain, "add_alias")
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			extraVars = vars
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute domain_management.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Adding domain: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("playbooks/domain_management.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Domain addition failed: %v", err)
			os.Exit(1)
		}

		// Add domain to configuration
		newDomain := models.Domain{
			Domain:     input.Domain,
			SSLEnabled: false,
		}

		stateMgr := state.NewManager(mgr)
		if err := stateMgr.AddDomainToSite(input.ServerName, input.SiteID, newDomain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		color.Green("\n✓ Domain '%s' added successfully", input.Domain)

		// Issue SSL if requested
		if input.IssueSSL {
			fmt.Println()
			color.Cyan("═══════════════════════════════════════════════════════")
			color.Cyan("  Issuing SSL certificate for: %s", input.Domain)
			color.Cyan("═══════════════════════════════════════════════════════")
			fmt.Println()

			sslVars := map[string]interface{}{
				"operation": "issue_ssl",
				"domain":    input.Domain,
			}
			if alias {
				vars, err := aliasVars(targetServer, input.SiteID, input.Domain, "issue_alias_ssl")
				if err != nil {
					color.Red("Error: %v", err)
					os.Exit(1)
				}
				sslVars = vars
			}

			sslResult, err := executor.ExecutePlaybookWithResult("playbooks/domain_management.yml", *targetServer, sslVars, cfg.GlobalVars)
			if err != nil {
				color.Red("\n✗ SSL certificate issuance failed: %v", err)
				fmt.Println("The domain has been added but SSL is not configured.")
				fmt.Println("You can issue SSL later with: press domain ssl")
				os.Exit(1)
			}

			// Update domain with SSL info
			now := time.Now()
			var expiresAt *time.Time

			// Try to parse actual expiry from Ansible output
			if sslResult.SSLInfo != nil && sslResult.SSLInfo.Expiry != "" {
				expiresAt = utils.ParseSSLExpiry(sslResult.SSLInfo.Expiry)
			}

			// Fallback to 90 days if parsing fails
			if expiresAt == nil {
				fallback := now.AddDate(0, 3, 0)
				expiresAt = &fallback
			}

			sslDomain := models.Domain{
				Domain:       input.Domain,
				SSLEnabled:   true,
				SSLIssuedAt:  &now,
				SSLExpiresAt: expiresAt,
			}

			if err := stateMgr.UpdateDomainSSL(input.ServerName, input.SiteID, input.Domain, sslDomain); err != nil {
				color.Red("Warning: Failed to update SSL status in configuration: %v", err)
			}

			color.Green("\n✓ SSL certificate issued successfully")
			fmt.Println()
			fmt.Printf("Domain URL:  https://%s\n", input.Domain)
			fmt.Printf("Expires:     %s\n", expiresAt.Format("2006-01-02"))
		} else {
			fmt.Println()
			fmt.Printf("Domain URL:  http://%s\n", input.Domain)
			fmt.Println()
			fmt.Println("To issue SSL later: press domain ssl")
		}
	},
}

// domainRemoveCmd represents the domain remove command
var domainRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete"},
	Short:   "Remove a domain from a site",
	Long: `Remove a domain from a WordPress site and its Nginx configuration.

Examples:
  # Interactive mode
  press domain remove

  # Non-interactive mode (for automation/AI agents)
  press domain remove --server myserver --site mysite --domain www.example.com --force`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		var input *prompt.DomainRemoveInput

		// Check for non-interactive mode
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")

		if serverName != "" && siteName != "" && domain != "" {
			// Non-interactive mode
			input = &prompt.DomainRemoveInput{
				ServerName: serverName,
				SiteID:     siteName,
				Domain:     domain,
			}
		} else if serverName != "" || siteName != "" || domain != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are all required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - get input from prompts
			var err error
			input, err = prompt.PromptDomainRemove(cfg.Servers)
			if err != nil {
				outputError(cmd, "Failed to get domain details", err)
				os.Exit(1)
			}
		}

		// The domain becomes root file paths and nginx config on the server.
		input.Domain = strings.ToLower(strings.TrimSpace(input.Domain))
		if err := utils.ValidateDomain(input.Domain); err != nil {
			outputError(cmd, "Invalid domain", fmt.Errorf("%q: %w", input.Domain, err))
			os.Exit(1)
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		// Only a domain the named site actually has: the playbook deletes the
		// domain's nginx configuration as root, whichever site it belongs to.
		if owner := utils.FindSiteByDomain(targetServer, input.Domain); owner == nil || owner.SiteID != input.SiteID {
			color.Red("Error: %s is not a domain of site '%s' on server '%s'", input.Domain, input.SiteID, input.ServerName)
			os.Exit(1)
		}

		// Final confirmation
		color.Yellow("\n⚠️  WARNING: This will remove:")
		fmt.Printf("  - Domain: %s\n", input.Domain)
		fmt.Printf("  - Nginx configuration\n")
		fmt.Printf("  - SSL certificate (if any)\n")
		fmt.Println()

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := huh.NewConfirm().
				Title("Remove this domain?").
				Affirmative("Yes, remove").
				Negative("Cancel").
				Value(&confirm).
				Run(); err != nil {
				os.Exit(1)
			}
			if !confirm {
				fmt.Println("Domain removal cancelled")
				return
			}
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"operation": "remove_domain",
			"domain":    input.Domain,
		}
		if alias, _ := cmd.Flags().GetBool("alias"); alias {
			vars, err := aliasVars(targetServer, input.SiteID, input.Domain, "remove_alias")
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			extraVars = vars
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute domain_management.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Removing domain: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("playbooks/domain_management.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Domain removal failed: %v", err)
			os.Exit(1)
		}

		// Remove domain from configuration
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.RemoveDomainFromSite(input.ServerName, input.SiteID, input.Domain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		color.Green("\n✓ Domain '%s' removed successfully", input.Domain)
	},
}

// domainSSLCmd represents the domain ssl command
var domainSSLCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Issue SSL certificate for a domain",
	Long: `Obtain a Let's Encrypt SSL certificate for a domain.

Examples:
  # Interactive mode
  press domain ssl

  # Non-interactive mode (for automation/AI agents)
  press domain ssl --server myserver --site mysite --domain www.example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		var input *prompt.DomainSSLInput

		// Check for non-interactive mode
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")

		if serverName != "" && siteName != "" && domain != "" {
			// Non-interactive mode
			input = &prompt.DomainSSLInput{
				ServerName: serverName,
				SiteID:     siteName,
				Domain:     domain,
			}
		} else if serverName != "" || siteName != "" || domain != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are all required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - get input from prompts
			var err error
			input, err = prompt.PromptDomainSSL(cfg.Servers)
			if err != nil {
				outputError(cmd, "Failed to get SSL details", err)
				os.Exit(1)
			}
		}

		input.Domain = strings.ToLower(strings.TrimSpace(input.Domain))
		if err := utils.ValidateDomain(input.Domain); err != nil {
			outputError(cmd, "Invalid domain", fmt.Errorf("%q: %w", input.Domain, err))
			os.Exit(1)
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"operation": "issue_ssl",
			"domain":    input.Domain,
		}
		if alias, _ := cmd.Flags().GetBool("alias"); alias {
			vars, err := aliasVars(targetServer, input.SiteID, input.Domain, "issue_alias_ssl")
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			extraVars = vars
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute domain_management.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Issuing SSL certificate for: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		result, err := executor.ExecutePlaybookWithResult("playbooks/domain_management.yml", *targetServer, extraVars, cfg.GlobalVars)
		if err != nil {
			color.Red("\n✗ SSL certificate issuance failed: %v", err)
			os.Exit(1)
		}

		// Update domain with SSL info
		now := time.Now()
		var expiresAt *time.Time

		// Try to parse actual expiry from Ansible output
		if result.SSLInfo != nil && result.SSLInfo.Expiry != "" {
			expiresAt = utils.ParseSSLExpiry(result.SSLInfo.Expiry)
		}

		// Fallback to 90 days if parsing fails
		if expiresAt == nil {
			fallback := now.AddDate(0, 3, 0)
			expiresAt = &fallback
		}

		sslDomain := models.Domain{
			Domain:       input.Domain,
			SSLEnabled:   true,
			SSLIssuedAt:  &now,
			SSLExpiresAt: expiresAt,
		}

		stateMgr := state.NewManager(mgr)
		if err := stateMgr.UpdateDomainSSL(input.ServerName, input.SiteID, input.Domain, sslDomain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ SSL certificate issued successfully!")
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Printf("Domain:      https://%s\n", input.Domain)
		fmt.Printf("Issued:      %s\n", now.Format("2006-01-02"))
		fmt.Printf("Expires:     %s\n", expiresAt.Format("2006-01-02"))
		fmt.Printf("Auto-renew:  Certbot will auto-renew before expiration\n")
	},
}

// primaryDomainForSite returns the primary domain of the site with the given
// ID on the server, or "" if not found.
func primaryDomainForSite(server *models.Server, siteID string) string {
	if server == nil {
		return ""
	}
	for i := range server.Sites {
		if server.Sites[i].SiteID == siteID {
			return server.Sites[i].PrimaryDomain
		}
	}
	return ""
}

// aliasVars builds the Ansible vars for managing an extra domain of a site.
// The site keeps living under /sites/<PrimaryDomain>; the alias only adds an
// nginx server block and a certificate.
func aliasVars(server *models.Server, siteID, domain, operation string) (map[string]interface{}, error) {
	for i := range server.Sites {
		site := server.Sites[i]
		if site.SiteID != siteID {
			continue
		}
		return map[string]interface{}{
			"operation":      operation,
			"domain":         domain,
			"site_id":        site.SiteID,
			"primary_domain": site.PrimaryDomain,
			"php_version":    site.PHPVersion,
		}, nil
	}
	return nil, fmt.Errorf("site '%s' not found on server '%s'", siteID, server.Name)
}

// domainPrimaryCmd makes a domain the address WordPress uses.
var domainPrimaryCmd = &cobra.Command{
	Use:   "primary",
	Short: "Make a domain the site's primary address",
	Long: `Switch WordPress's home/siteurl (and every stored URL) to the given domain.
The domain must already be added to the site.

Example:
  press domain primary --server myserver --site mysite --domain www.example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")
		domain, _ := cmd.Flags().GetString("domain")
		if serverName == "" || siteName == "" || domain == "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("--server, --site, and --domain are required"))
			os.Exit(1)
		}
		domain = strings.ToLower(strings.TrimSpace(domain))
		if err := utils.ValidateDomain(domain); err != nil {
			outputError(cmd, "Invalid domain", fmt.Errorf("%q: %w", domain, err))
			os.Exit(1)
		}

		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == serverName {
				targetServer = &cfg.Servers[i]
				break
			}
		}
		if targetServer == nil {
			color.Red("Error: Server '%s' not found", serverName)
			os.Exit(1)
		}

		vars, err := aliasVars(targetServer, siteName, domain, "set_primary")
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		if err := executor.ExecutePlaybook("playbooks/domain_management.yml", *targetServer, vars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Switching the primary domain failed: %v", err)
			os.Exit(1)
		}

		if err := state.NewManager(mgr).SetSitePublicDomain(serverName, siteName, domain); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		color.Green("\n✓ '%s' is now the primary domain", domain)
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domainAddCmd)
	domainCmd.AddCommand(domainRemoveCmd)
	domainCmd.AddCommand(domainSSLCmd)
	domainCmd.AddCommand(domainPrimaryCmd)

	domainPrimaryCmd.Flags().String("server", "", "Server name")
	domainPrimaryCmd.Flags().String("site", "", "Site ID")
	domainPrimaryCmd.Flags().String("domain", "", "Domain to make primary")

	// domain add flags (non-interactive mode)
	domainAddCmd.Flags().String("server", "", "Server name")
	domainAddCmd.Flags().String("site", "", "Site ID")
	domainAddCmd.Flags().String("domain", "", "Domain to add")
	domainAddCmd.Flags().Bool("ssl", false, "Issue SSL certificate for the domain")
	domainAddCmd.Flags().Bool("json", false, "Output in JSON format")
	domainAddCmd.Flags().Bool("alias", false, "The domain is an extra domain of an existing site (served from the site's own directory)")

	// domain remove flags
	domainRemoveCmd.Flags().String("server", "", "Server name")
	domainRemoveCmd.Flags().String("site", "", "Site ID")
	domainRemoveCmd.Flags().String("domain", "", "Domain to remove")
	domainRemoveCmd.Flags().BoolP("force", "f", false, "Force removal without confirmation")
	domainRemoveCmd.Flags().Bool("json", false, "Output in JSON format")
	domainRemoveCmd.Flags().Bool("alias", false, "The domain is an extra domain of an existing site (served from the site's own directory)")

	// domain ssl flags (non-interactive mode)
	domainSSLCmd.Flags().String("server", "", "Server name")
	domainSSLCmd.Flags().String("site", "", "Site ID")
	domainSSLCmd.Flags().String("domain", "", "Domain to issue SSL for")
	domainSSLCmd.Flags().Bool("json", false, "Output in JSON format")
	domainSSLCmd.Flags().Bool("alias", false, "The domain is an extra domain of an existing site (served from the site's own directory)")
}
