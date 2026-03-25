package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"

	// Global flags
	Verbose bool
	DryRun  bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "press",
	Short: "pressctl - Ansible wrapper for WordPress hosting management",
	Long: `pressctl is a CLI tool that simplifies WordPress hosting management
by wrapping Ansible playbooks with an intuitive, interactive interface.

Manage servers, sites, and domains with ease while maintaining full
visibility into your infrastructure state via ~/.pressctl/pressctl.yaml

Examples:
  # Initialize configuration
  press init

  # Add and provision a new server
  press server provision

  # Create a WordPress site
  press site create

  # Add a domain with SSL
  press domain add

  # List all servers
  press server list --json`,
	Run: func(cmd *cobra.Command, args []string) {
		type action struct {
			label string
			run   func()
		}

		actions := []action{
			{"Provision a new server", func() { serverProvisionCmd.Run(serverProvisionCmd, []string{}) }},
			{"Create a WordPress site", func() { siteCreateCmd.Run(siteCreateCmd, []string{}) }},
			{"Delete a site", func() { siteDeleteCmd.Run(siteDeleteCmd, []string{}) }},
			{"Add a domain to a site", func() { domainAddCmd.Run(domainAddCmd, []string{}) }},
			{"Issue / renew SSL for a domain", func() { domainSSLCmd.Run(domainSSLCmd, []string{}) }},
			{"Check server health", func() { serverHealthCheckCmd.Run(serverHealthCheckCmd, []string{}) }},
			{"List servers", func() { serverListCmd.Run(serverListCmd, []string{}) }},
			{"List sites", func() { siteListCmd.Run(siteListCmd, []string{}) }},
		}

		fmt.Println()
		color.Cyan("  pressctl — Common Actions")
		fmt.Println()
		for i, a := range actions {
			fmt.Printf("  %s  %s\n", color.YellowString("%d.", i+1), a.label)
		}
		fmt.Println()

		var input string
		if err := survey.AskOne(&survey.Input{Message: "Enter number:"}, &input); err != nil {
			os.Exit(1)
		}

		n, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || n < 1 || n > len(actions) {
			color.Red("Invalid selection. Please enter a number between 1 and %d.", len(actions))
			os.Exit(1)
		}

		fmt.Println()
		actions[n-1].run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "Show what would be done without making changes")
}
