package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/pressctl/cli/internal/ansible"
	"github.com/pressctl/cli/internal/state"
	"github.com/pressctl/cli/pkg/models"
	"github.com/spf13/cobra"
)

// serverPhpVersionCmd represents the server php-version command group
var serverPhpVersionCmd = &cobra.Command{
	Use:     "php-version",
	Aliases: []string{"php"},
	Short:   "Manage PHP versions on a server",
	Long: `Manage PHP versions installed on a provisioned server.

Install additional PHP versions so sites can be migrated to them, or list
the versions currently installed.

Examples:
  press server php-version add
  press server php-version add --server myserver --version 8.4
  press server php-version list`,
}

// serverPhpVersionAddCmd installs an additional PHP version on a server
var serverPhpVersionAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"install"},
	Short:   "Install an additional PHP version on a server",
	Long: `Install an additional PHP version (and its extensions) on a provisioned
server so sites can be migrated to it later with 'press site php-version'.

Examples:
  # Interactive mode
  press server php-version add

  # Non-interactive mode (for automation/AI agents)
  press server php-version add --server myserver --version 8.4`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		serverName, _ := cmd.Flags().GetString("server")
		version, _ := cmd.Flags().GetString("version")

		// Partial flags are an error
		if (serverName == "") != (version == "") {
			color.Red("Error: --server and --version are both required for non-interactive mode")
			os.Exit(1)
		}

		// Find the target server (select interactively if not provided)
		var targetServer *models.Server
		if serverName == "" {
			options := make([]string, 0)
			for i := range cfg.Servers {
				options = append(options, cfg.Servers[i].Name)
			}
			if len(options) == 0 {
				color.Red("Error: No servers configured")
				fmt.Println("Add a server first: press server add")
				os.Exit(1)
			}

			opts := make([]huh.Option[int], len(options))
			for i, o := range options {
				opts[i] = huh.NewOption(o, i)
			}
			var selected int
			if err := huh.NewSelect[int]().
				Title("Select a server").
				Options(opts...).
				Value(&selected).
				Run(); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			serverName = options[selected]
		}

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

		if targetServer.Status != "provisioned" {
			color.Red("Error: Server '%s' is not provisioned", serverName)
			fmt.Println("Provision the server first: press server provision", serverName)
			os.Exit(1)
		}

		// Choose the PHP version to install (interactively if not provided)
		if version == "" {
			available := make([]string, 0)
			for _, v := range models.SupportedPHPVersions {
				if !targetServer.HasPHPVersion(v) {
					available = append(available, v)
				}
			}
			if len(available) == 0 {
				fmt.Printf("All supported PHP versions are already installed on '%s': %s\n",
					serverName, strings.Join(targetServer.InstalledPHPVersions(), ", "))
				return
			}

			opts := make([]huh.Option[string], len(available))
			for i, o := range available {
				opts[i] = huh.NewOption(o, o)
			}
			if err := huh.NewSelect[string]().
				Title("PHP version to install").
				Description(fmt.Sprintf("Currently installed: %s", strings.Join(targetServer.InstalledPHPVersions(), ", "))).
				Options(opts...).
				Value(&version).
				Run(); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
		}

		// Validate the target PHP version
		if !models.IsValidPHPVersion(version) {
			color.Red("Error: Unsupported PHP version '%s'", version)
			fmt.Printf("Supported versions: %s\n", strings.Join(models.SupportedPHPVersions, ", "))
			os.Exit(1)
		}

		if targetServer.HasPHPVersion(version) {
			color.Yellow("PHP %s is already installed on server '%s'", version, serverName)
			fmt.Printf("Installed versions: %s\n", strings.Join(targetServer.InstalledPHPVersions(), ", "))
			return
		}

		// Confirm installation
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Install PHP %s on '%s'?", version, serverName)).
				Affirmative("Yes, install").
				Negative("Cancel").
				Value(&confirm).
				Run(); err != nil {
				os.Exit(1)
			}
			if !confirm {
				fmt.Println("Installation cancelled")
				return
			}
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"php_version": version,
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute php_version_add playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Installing PHP %s on server: %s", version, serverName)
		color.Cyan("  Estimated time: 1-3 minutes")
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("playbooks/php_version_add.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ PHP %s installation failed: %v", version, err)
			os.Exit(1)
		}

		// Update server configuration
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.AddServerPHPVersion(serverName, version); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		fmt.Println()
		color.Green("✓ PHP %s installed on server '%s'", version, serverName)
		fmt.Println()
		fmt.Println("Next step: point a site at this version with: press site php-version")
	},
}

// serverPhpVersionListCmd lists the PHP versions installed on a server
var serverPhpVersionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List PHP versions installed on a server",
	Long: `Display the PHP versions installed on a server.

Examples:
  # Interactive mode
  press server php-version list

  # List versions on a specific server
  press server php-version list --server myserver`,
	Run: func(cmd *cobra.Command, args []string) {
		_, cfg := ensureConfig()

		serverName, _ := cmd.Flags().GetString("server")

		// Select server interactively if not provided
		if serverName == "" {
			options := make([]string, 0)
			for i := range cfg.Servers {
				options = append(options, cfg.Servers[i].Name)
			}
			if len(options) == 0 {
				color.Red("Error: No servers configured")
				fmt.Println("Add a server first: press server add")
				os.Exit(1)
			}

			opts := make([]huh.Option[int], len(options))
			for i, o := range options {
				opts[i] = huh.NewOption(o, i)
			}
			var selected int
			if err := huh.NewSelect[int]().
				Title("Select a server").
				Options(opts...).
				Value(&selected).
				Run(); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
			serverName = options[selected]
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

		versions := targetServer.InstalledPHPVersions()

		// Check for JSON output
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			type serverPHPOutput struct {
				Server   string   `json:"server"`
				Versions []string `json:"php_versions"`
			}
			output, err := json.MarshalIndent(serverPHPOutput{Server: serverName, Versions: versions}, "", "  ")
			if err != nil {
				color.Red("Error: Failed to marshal JSON: %v", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
			return
		}

		if len(versions) == 0 {
			fmt.Printf("No PHP versions tracked for server '%s'.\n", serverName)
			fmt.Println("The default version is set at provisioning time (server php_versions).")
			return
		}

		fmt.Printf("\nPHP versions on server '%s':\n\n", serverName)
		for i, v := range versions {
			defaultMark := ""
			if v == targetServer.PHPVersion {
				defaultMark = " (default)"
			}
			fmt.Printf("  %d. PHP %s%s\n", i+1, v, defaultMark)
		}
		fmt.Println()
	},
}

func init() {
	serverCmd.AddCommand(serverPhpVersionCmd)
	serverPhpVersionCmd.AddCommand(serverPhpVersionAddCmd)
	serverPhpVersionCmd.AddCommand(serverPhpVersionListCmd)

	// server php-version add flags
	serverPhpVersionAddCmd.Flags().String("server", "", "Server name")
	serverPhpVersionAddCmd.Flags().String("version", "", "PHP version to install (8.1, 8.2, 8.3, 8.4, 8.5)")
	serverPhpVersionAddCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	serverPhpVersionAddCmd.Flags().Bool("json", false, "Output in JSON format")

	// server php-version list flags
	serverPhpVersionListCmd.Flags().String("server", "", "Server name")
	serverPhpVersionListCmd.Flags().Bool("json", false, "Output in JSON format")
}
