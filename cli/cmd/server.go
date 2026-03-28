package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/pressctl/cli/internal/ansible"
	"github.com/pressctl/cli/internal/prompt"
	"github.com/pressctl/cli/internal/state"
	"github.com/pressctl/cli/internal/utils"
	"github.com/pressctl/cli/pkg/models"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	Long: `Add, list, remove, and provision servers.

The 'provision' command provides an interactive way to add and provision a new server in one step.`,
}

// serverAddCmd represents the server add command
var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new server without provisioning",
	Long: `Add a new server to the configuration without provisioning it.

For most users, use 'press server provision' instead, which adds and provisions in one step.

Examples:
  # Interactive mode
  press server add

  # Non-interactive mode (for automation/AI agents)
  press server add --name myserver --ip 1.2.3.4 --ssh-user root`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		var input *prompt.ServerInput

		// Check for non-interactive mode
		name, _ := cmd.Flags().GetString("name")
		ip, _ := cmd.Flags().GetString("ip")

		if name != "" && ip != "" {
			// Non-interactive mode
			sshUser, _ := cmd.Flags().GetString("ssh-user")
			sshPort, _ := cmd.Flags().GetInt("ssh-port")

			input = &prompt.ServerInput{
				Name:     name,
				Hostname: ip,
				IP:       ip,
				SSHUser:  sshUser,
				SSHPort:  sshPort,
			}

			if input.SSHUser == "" {
				input.SSHUser = "root"
			}
			if input.SSHPort == 0 {
				input.SSHPort = 22
			}
		} else if name != "" || ip != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("both --name and --ip are required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode - prompt for server details
			var err error
			input, err = prompt.PromptServerAdd()
			if err != nil {
				outputError(cmd, "Failed to get server details", err)
				os.Exit(1)
			}
		}

		// Check for duplicate server name
		for _, server := range cfg.Servers {
			if server.Name == input.Name {
				outputError(cmd, "Server already exists", fmt.Errorf("server with name '%s' already exists", input.Name))
				os.Exit(1)
			}
		}

		// Add server to config
		newServer := input.ToServer()
		cfg.Servers = append(cfg.Servers, newServer)

		// Save config
		if err := mgr.Save(cfg); err != nil {
			outputError(cmd, "Failed to save configuration", err)
			os.Exit(1)
		}

		outputSuccess(cmd, "server_added", map[string]interface{}{
			"name":   input.Name,
			"ip":     input.IP,
			"status": "unprovisioned",
		})

		if !isJSONOutput(cmd) {
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Printf("  Provision the server: press server provision %s\n", input.Name)
		}
	},
}

// serverListCmd represents the server list command
var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Long:  `Display all servers in the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, cfg := ensureConfig()

		// Check for JSON output
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(cfg.Servers, "", "  ")
			if err != nil {
				color.Red("Error: Failed to marshal JSON: %v", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
			return
		}

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			fmt.Println("Add and provision a server with: press server provision")
			return
		}

		fmt.Printf("\nServers (%d total):\n\n", len(cfg.Servers))

		// Prepare table data
		headers := []string{"NAME", "HOSTNAME", "IP", "SSH USER", "STATUS", "SITES"}
		colWidths := []int{18, 28, 15, 12, 15, 6}
		rows := make([][]string, 0)

		for _, server := range cfg.Servers {
			statusStr := ""
			switch server.Status {
			case "provisioned":
				statusStr = color.GreenString(server.Status)
			case "unprovisioned":
				statusStr = color.YellowString(server.Status)
			case "error":
				statusStr = color.RedString(server.Status)
			default:
				statusStr = server.Status
			}

			row := []string{
				server.Name,
				server.Hostname,
				server.IP,
				server.SSH.User,
				statusStr,
				fmt.Sprintf("%d", len(server.Sites)),
			}
			rows = append(rows, row)
		}

		utils.PrintTableWithBorders(headers, rows, colWidths)
	},
}

// serverRemoveCmd represents the server remove command
var serverRemoveCmd = &cobra.Command{
	Use:     "remove [name]",
	Aliases: []string{"delete"},
	Short:   "Remove a server from inventory",
	Long: `Remove a server from the pressctl inventory.

If no name is provided, you will be prompted to select a server to remove.

Note: This only removes the server from the pressctl inventory. The actual server
and its resources will still exist in your cloud provider. You must manually
delete the server from your cloud provider (AWS, DigitalOcean, etc.) if needed.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			return
		}

		var serverName string

		// Interactive mode: no server name provided
		if len(args) == 0 {
			// Build options list
			options := make([]string, len(cfg.Servers))
			for i, server := range cfg.Servers {
				siteCount := len(server.Sites)
				siteLabel := "sites"
				if siteCount == 1 {
					siteLabel = "site"
				}
				options[i] = fmt.Sprintf("%s (%s) - %d %s", server.Name, server.IP, siteCount, siteLabel)
			}

			opts := make([]huh.Option[int], len(options))
			for i, o := range options {
				opts[i] = huh.NewOption(o, i)
			}
			var selected int
			if err := huh.NewSelect[int]().
				Title("Select a server to remove").
				Options(opts...).
				Value(&selected).
				Run(); err != nil {
				os.Exit(1)
			}
			serverName = cfg.Servers[selected].Name
		} else {
			serverName = args[0]
		}

		// Find and remove server
		found := false
		newServers := make([]models.Server, 0)
		var removedServer models.Server

		for _, server := range cfg.Servers {
			if server.Name == serverName {
				found = true
				removedServer = server
			} else {
				newServers = append(newServers, server)
			}
		}

		if !found {
			color.Red("Error: Server '%s' not found", serverName)
			os.Exit(1)
		}

		// Show warning about cloud provider
		fmt.Println()
		color.Yellow("Warning: This will remove '%s' from the pressctl inventory only.", serverName)
		fmt.Println("The server will still exist in your cloud provider.")
		fmt.Println("You must manually delete it from your cloud provider if needed.")
		fmt.Println()

		// Warn if server has sites
		if len(removedServer.Sites) > 0 {
			color.Yellow("This server has %d site(s) that will also be removed from the inventory.", len(removedServer.Sites))
			fmt.Println()
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Remove server '%s' from inventory?", serverName)).
				Affirmative("Yes, remove").
				Negative("Cancel").
				Value(&confirm).
				Run(); err != nil {
				os.Exit(1)
			}
			if !confirm {
				fmt.Println("Server removal cancelled")
				return
			}
		}

		cfg.Servers = newServers

		// Save config
		if err := mgr.Save(cfg); err != nil {
			color.Red("Error: Failed to save configuration: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Server '%s' removed from inventory", serverName)
	},
}

// serverProvisionCmd represents the server provision command
var serverProvisionCmd = &cobra.Command{
	Use:   "provision [name]",
	Short: "Provision a server",
	Long: `Run the provision.yml playbook to set up a server with Nginx, PHP, MariaDB, and security hardening.

If no name is provided, you will be prompted to add a new server and provision it immediately.
If a name is provided, the existing server will be provisioned.

Examples:
  # Interactive mode - add and provision new server
  press server provision

  # Provision existing server by name
  press server provision myserver

  # Non-interactive mode - add and provision new server (for automation/AI agents)
  press server provision --name myserver --ip 1.2.3.4 --force`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		var targetServer *models.Server
		var serverName string

		// Check for non-interactive mode via flags
		flagName, _ := cmd.Flags().GetString("name")
		flagIP, _ := cmd.Flags().GetString("ip")

		if len(args) > 0 {
			// Provision existing server by name argument
			serverName = args[0]

			// Find the server
			for i := range cfg.Servers {
				if cfg.Servers[i].Name == serverName {
					targetServer = &cfg.Servers[i]
					break
				}
			}

			if targetServer == nil {
				outputError(cmd, "Server not found", fmt.Errorf("server '%s' not found. Run 'press server list' to see available servers", serverName))
				os.Exit(1)
			}
		} else if flagName != "" && flagIP != "" {
			// Non-interactive mode: create new server from flags
			sshUser, _ := cmd.Flags().GetString("ssh-user")
			sshPort, _ := cmd.Flags().GetInt("ssh-port")
			flagPHP, _ := cmd.Flags().GetString("php-version")

			if flagPHP == "" {
				flagPHP = models.DefaultPHPVersion
			}

			// Check for duplicate server name
			for _, server := range cfg.Servers {
				if server.Name == flagName {
					outputError(cmd, "Server already exists", fmt.Errorf("server with name '%s' already exists", flagName))
					os.Exit(1)
				}
			}

			// Create new server
			newServer := models.Server{
				Name:     flagName,
				Hostname: flagIP,
				IP:       flagIP,
				SSH: models.SSHConfig{
					User: sshUser,
					Port: sshPort,
				},
				PHPVersion: flagPHP,
				Status:     "unprovisioned",
				Sites:      []models.Site{},
			}

			cfg.Servers = append(cfg.Servers, newServer)

			// Save config
			if err := mgr.Save(cfg); err != nil {
				outputError(cmd, "Failed to save configuration", err)
				os.Exit(1)
			}

			outputInfo(cmd, "✓ Server '%s' added to configuration\n\n", flagName)

			serverName = flagName
			targetServer = &cfg.Servers[len(cfg.Servers)-1]
		} else if flagName != "" || flagIP != "" {
			outputError(cmd, "Incomplete flags", fmt.Errorf("both --name and --ip are required for non-interactive mode"))
			os.Exit(1)
		} else {
			// Interactive mode: prompt for server details
			input, err := prompt.PromptServerAdd()
			if err != nil {
				outputError(cmd, "Failed to get server details", err)
				os.Exit(1)
			}

			// Check for duplicate server name
			for _, server := range cfg.Servers {
				if server.Name == input.Name {
					outputError(cmd, "Server already exists", fmt.Errorf("server with name '%s' already exists", input.Name))
					os.Exit(1)
				}
			}

			// Add server to config
			newServer := input.ToServer()
			cfg.Servers = append(cfg.Servers, newServer)

			// Save config
			if err := mgr.Save(cfg); err != nil {
				outputError(cmd, "Failed to save configuration", err)
				os.Exit(1)
			}

			color.Green("✓ Server '%s' added to configuration", input.Name)
			fmt.Println()

			// Set target server for provisioning
			serverName = input.Name
			targetServer = &cfg.Servers[len(cfg.Servers)-1]
		}

		// Check if already provisioned
		if targetServer.Status == "provisioned" {
			color.Yellow("Warning: Server '%s' is already marked as provisioned", serverName)

			skipCheck, _ := cmd.Flags().GetBool("skip-check")
			if !skipCheck {
				var confirm bool
				if err := huh.NewConfirm().
					Title("Provision again anyway?").
					Affirmative("Yes, reprovision").
					Negative("Cancel").
					Value(&confirm).
					Run(); err != nil {
					os.Exit(1)
				}
				if !confirm {
					fmt.Println("Provisioning cancelled")
					return
				}
			}
		}

		// Pre-flight SSH check
		skipSSH, _ := cmd.Flags().GetBool("skip-ssh-check")
		if !skipSSH {
			if Debug {
				fmt.Println("[debug] SSH config:")
				fmt.Printf("  user:     %s\n", targetServer.SSH.User)
				fmt.Printf("  host:     %s:%d\n", targetServer.IP, targetServer.SSH.Port)
				fmt.Println()
			}
			fmt.Println("Checking SSH connectivity...")
			if err := utils.TestSSHConnection(*targetServer); err != nil {
				color.Red("✗ SSH connectivity check failed: %v", err)
				fmt.Println()
				fmt.Println("Please verify:")
				fmt.Println("  1. Server is reachable")
				fmt.Println("  2. SSH agent is running with your key loaded (ssh-add -l)")
				fmt.Println("  3. SSH user has access to the server")
				fmt.Println()
				fmt.Println("Use --skip-ssh-check to bypass this check (not recommended)")
				os.Exit(1)
			}
			color.Green("✓ SSH connectivity check passed")
			fmt.Println()
		}

		// Pre-flight port conflict check
		skipPortCheck, _ := cmd.Flags().GetBool("skip-port-check")
		if !skipPortCheck {
			fmt.Println("Checking for port conflicts...")
			conflicts, portErr := utils.CheckPortConflicts(*targetServer)
			if portErr != nil {
				color.Yellow("⚠ Port conflict check skipped: %v", portErr)
				fmt.Println()
			} else if len(conflicts) > 0 {
				color.Red("✗ Port conflicts detected on %s:", targetServer.IP)
				for _, c := range conflicts {
					color.Red("  • %s", c)
				}
				fmt.Println()
				fmt.Println("These services will conflict with provisioning:")
				fmt.Println("  - A web server on ports 80/443 will prevent Nginx from binding")
				fmt.Println("  - MySQL/MariaDB on port 3306 may conflict with the MariaDB install")
				fmt.Println()
				fmt.Println("Please stop conflicting services before provisioning.")
				fmt.Println("Use --skip-port-check to bypass this check (not recommended).")
				fmt.Println()

				force, _ := cmd.Flags().GetBool("force")
				if !force {
					var proceed bool
					if err := huh.NewConfirm().
						Title("Proceed anyway?").
						Affirmative("Yes, proceed").
						Negative("Cancel").
						Value(&proceed).
						Run(); err != nil {
						os.Exit(1)
					}
					if !proceed {
						fmt.Println("Provisioning cancelled")
						return
					}
				}
			} else {
				color.Green("✓ No port conflicts detected")
				fmt.Println()
			}
		}

		// Confirm provisioning
		phpVersion := targetServer.PHPVersion
		if phpVersion == "" {
			phpVersion = models.DefaultPHPVersion
		}
		color.Cyan("About to provision server: %s (%s)", targetServer.Name, targetServer.IP)
		fmt.Println("This will:")
		fmt.Printf("  - Install Nginx, PHP %s, MariaDB\n", phpVersion)
		fmt.Println("  - Configure security (UFW, Fail2ban, SSH hardening)")
		fmt.Println("  - Set up Certbot for SSL certificates")
		fmt.Println("  - Create pressctl user and environment")
		fmt.Println()

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := huh.NewConfirm().
				Title("Continue with provisioning?").
				Affirmative("Yes, provision").
				Negative("Cancel").
				Value(&confirm).
				Run(); err != nil {
				os.Exit(1)
			}
			if !confirm {
				fmt.Println("Provisioning cancelled")
				return
			}
		}

		// Generate MySQL password for this server if not already set
		mysqlPassword := targetServer.Credentials.MySQLWordmonbotPassword
		if mysqlPassword == "" {
			mysqlPassword = prompt.GenerateSecurePassword(24)
			targetServer.Credentials.MySQLWordmonbotPassword = mysqlPassword

			// Update server in config with the new password
			for i := range cfg.Servers {
				if cfg.Servers[i].Name == serverName {
					cfg.Servers[i].Credentials.MySQLWordmonbotPassword = mysqlPassword
					break
				}
			}
			if err := mgr.Save(cfg); err != nil {
				outputError(cmd, "Failed to save MySQL password to config", err)
				os.Exit(1)
			}
		}

		// Build provision vars
		provisionVars := make(map[string]interface{})
		for k, v := range cfg.GlobalVars {
			provisionVars[k] = v
		}
		provisionVars["mysql_pressctlbot_password"] = mysqlPassword
		provisionVars["php_version"] = phpVersion

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute provision.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Starting provisioning: %s", serverName)
		color.Cyan("  Estimated time: 5-10 minutes")
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("provision.yml", *targetServer, nil, provisionVars); err != nil {
			color.Red("\n✗ Provisioning failed: %v", err)

			// Mark server as error
			stateMgr := state.NewManager(mgr)
			stateMgr.MarkServerError(serverName)

			os.Exit(1)
		}

		// Update server status to provisioned
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.MarkServerProvisioned(serverName); err != nil {
			color.Red("Warning: Failed to update server status: %v", err)
		}

		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ Server '%s' provisioned successfully!", serverName)
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Server credentials:")
		fmt.Printf("  MySQL pressctlbot password: %s\n", mysqlPassword)
		fmt.Println()
		color.Yellow("  Save this password! It's stored in your config file.")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  Create a WordPress site: press site create")
	},
}

// serverHealthCheckCmd represents the server health-check command
var serverHealthCheckCmd = &cobra.Command{
	Use:     "health-check [name]",
	Aliases: []string{"check", "ping"},
	Short:   "Check server connectivity and health",
	Long: `Test SSH connectivity and verify services are running on a server.

Examples:
  # Check a specific server
  press server health-check myserver

  # Interactively select a server to check
  press server health-check`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, cfg := ensureConfig()

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			return
		}

		var serverName string

		if len(args) == 0 {
			// Interactive mode
			options := make([]string, len(cfg.Servers))
			for i, server := range cfg.Servers {
				options[i] = fmt.Sprintf("%s (%s) - %s", server.Name, server.IP, server.Status)
			}

			opts := make([]huh.Option[int], len(options))
			for i, o := range options {
				opts[i] = huh.NewOption(o, i)
			}
			var selected int
			if err := huh.NewSelect[int]().
				Title("Select a server to check").
				Options(opts...).
				Value(&selected).
				Run(); err != nil {
				os.Exit(1)
			}
			serverName = cfg.Servers[selected].Name
		} else {
			serverName = args[0]
		}

		// Find server
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

		fmt.Printf("\nChecking server: %s (%s)\n\n", targetServer.Name, targetServer.IP)

		// Test SSH connectivity
		fmt.Print("SSH connectivity... ")
		if err := utils.TestSSHConnection(*targetServer); err != nil {
			color.Red("FAILED")
			color.Red("  %v", err)
			os.Exit(1)
		}
		color.Green("OK")

		fmt.Println()
		color.Green("✓ Server '%s' is healthy", serverName)
	},
}

// serverUpdateCmd represents the server update command
var serverUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update server configuration",
	Long: `Update the configuration for an existing server.

Examples:
  # Update a specific server
  press server update myserver

  # Interactively select a server to update
  press server update`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr, cfg := ensureConfig()

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			return
		}

		var serverName string

		if len(args) == 0 {
			// Interactive mode
			options := make([]string, len(cfg.Servers))
			for i, server := range cfg.Servers {
				options[i] = fmt.Sprintf("%s (%s)", server.Name, server.IP)
			}

			opts := make([]huh.Option[int], len(options))
			for i, o := range options {
				opts[i] = huh.NewOption(o, i)
			}
			var selected int
			if err := huh.NewSelect[int]().
				Title("Select a server to update").
				Options(opts...).
				Value(&selected).
				Run(); err != nil {
				os.Exit(1)
			}
			serverName = cfg.Servers[selected].Name
		} else {
			serverName = args[0]
		}

		// Find server index
		var serverIndex int = -1
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == serverName {
				serverIndex = i
				break
			}
		}

		if serverIndex == -1 {
			color.Red("Error: Server '%s' not found", serverName)
			os.Exit(1)
		}

		server := &cfg.Servers[serverIndex]

		fmt.Printf("\nUpdating server: %s\n", server.Name)
		fmt.Println("Leave blank to keep current value.")

		// Update fields interactively
		var newName = server.Name
		var newIP = server.IP
		var newSSHUser = server.SSH.User
		var newSSHPort = fmt.Sprintf("%d", server.SSH.Port)

		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Server name").
					Value(&newName),
				huh.NewInput().
					Title("IP address").
					Value(&newIP).
					Validate(utils.ValidateIP),
				huh.NewInput().
					Title("SSH user").
					Value(&newSSHUser),
				huh.NewInput().
					Title("SSH port").
					Value(&newSSHPort).
					Validate(utils.ValidatePort),
			),
		).Run(); err != nil {
			os.Exit(1)
		}

		// Check for duplicate name
		if newName != server.Name {
			for _, s := range cfg.Servers {
				if s.Name == newName {
					color.Red("Error: Server with name '%s' already exists", newName)
					os.Exit(1)
				}
			}
		}

		var port int
		fmt.Sscanf(newSSHPort, "%d", &port)
		if port == 0 {
			port = 22
		}

		// Show summary and confirm
		fmt.Println()
		fmt.Println("Updated configuration:")
		fmt.Printf("  Name:     %s\n", newName)
		fmt.Printf("  IP:       %s\n", newIP)
		fmt.Printf("  SSH User: %s\n", newSSHUser)
		fmt.Printf("  SSH Port: %d\n", port)
		fmt.Println()

		var confirm bool
		if err := huh.NewConfirm().
			Title("Save changes?").
			Affirmative("Yes, save").
			Negative("Cancel").
			Value(&confirm).
			Run(); err != nil {
			os.Exit(1)
		}
		if !confirm {
			fmt.Println("Update cancelled")
			return
		}

		// Apply changes
		server.Name = newName
		server.Hostname = newIP
		server.IP = newIP
		server.SSH.User = newSSHUser
		server.SSH.Port = port

		// Save config
		if err := mgr.Save(cfg); err != nil {
			color.Red("Error: Failed to save configuration: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Server '%s' updated successfully", newName)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverAddCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverRemoveCmd)
	serverCmd.AddCommand(serverProvisionCmd)
	serverCmd.AddCommand(serverHealthCheckCmd)
	serverCmd.AddCommand(serverUpdateCmd)

	// server add flags (non-interactive mode)
	serverAddCmd.Flags().String("name", "", "Server name")
	serverAddCmd.Flags().String("ip", "", "Server IP address")
	serverAddCmd.Flags().String("ssh-user", "root", "SSH user")
	serverAddCmd.Flags().Int("ssh-port", 22, "SSH port")
	serverAddCmd.Flags().Bool("json", false, "Output in JSON format")

	// server list flags
	serverListCmd.Flags().Bool("json", false, "Output in JSON format")

	// server remove flags
	serverRemoveCmd.Flags().BoolP("force", "f", false, "Force removal without confirmation")
	serverRemoveCmd.Flags().Bool("json", false, "Output in JSON format")

	// server provision flags
	serverProvisionCmd.Flags().String("name", "", "Server name (for non-interactive mode)")
	serverProvisionCmd.Flags().String("ip", "", "Server IP address")
	serverProvisionCmd.Flags().String("ssh-user", "root", "SSH user")
	serverProvisionCmd.Flags().Int("ssh-port", 22, "SSH port")
	serverProvisionCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	serverProvisionCmd.Flags().Bool("skip-ssh-check", false, "Skip SSH connectivity check")
	serverProvisionCmd.Flags().Bool("skip-port-check", false, "Skip port conflict check")
	serverProvisionCmd.Flags().Bool("skip-check", false, "Skip already-provisioned check")
	serverProvisionCmd.Flags().String("php-version", "", "PHP version to install (8.1, 8.2, 8.3, 8.4)")
	serverProvisionCmd.Flags().Bool("json", false, "Output in JSON format")

	// server health-check flags
	serverHealthCheckCmd.Flags().Bool("json", false, "Output in JSON format")

	// server update flags
	serverUpdateCmd.Flags().String("name", "", "New server name")
	serverUpdateCmd.Flags().String("ip", "", "New IP address")
	serverUpdateCmd.Flags().String("ssh-user", "", "New SSH user")
	serverUpdateCmd.Flags().Int("ssh-port", 0, "New SSH port")
	serverUpdateCmd.Flags().Bool("json", false, "Output in JSON format")
}
