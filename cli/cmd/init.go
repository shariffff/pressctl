package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/pressctl/cli/internal/config"
	"github.com/pressctl/cli/internal/installer"
	"github.com/pressctl/cli/internal/prompt"
)

// ensureConfig loads the config, auto-initializing on first use if needed.
// Every command that needs config should call this instead of manually
// checking ConfigExists().
func ensureConfig() (*config.Manager, *config.Config) {
	mgr, err := config.NewManager()
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	if mgr.ConfigExists() {
		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}
		if migrated, err := config.MigrateIfNeeded(mgr, cfg); err != nil {
			color.Yellow("Warning: Config migration failed: %v", err)
		} else if migrated {
			color.Cyan("→ Config migrated to schema v%s", config.SchemaVersion)
		}
		return mgr, cfg
	}

	// Auto-initialize on first use
	fmt.Println()
	color.Cyan("First run — initializing pressctl...")
	fmt.Println()

	// Copy Ansible playbooks if needed
	if !installer.IsInitialized() {
		fmt.Print("→ Copying Ansible playbooks... ")
		if err := installer.Initialize(); err != nil {
			color.Red("✗")
			color.Red("\nError: %v", err)
			color.Red("Could not find Ansible playbooks. Make sure pressctl is installed correctly.")
			os.Exit(1)
		}
		color.Green("✓")
	}

	// Auto-detect SSH public key
	initInput, err := prompt.PromptInitSetup()
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	// Create config
	fmt.Print("→ Creating configuration... ")
	cfg := config.DefaultConfig()
	cfg.Ansible.Path = installer.GetAnsibleDir()
	cfg.GlobalVars["pressctl_ssh_key"] = initInput.SSHPublicKey

	if err := mgr.Save(cfg); err != nil {
		color.Red("✗")
		color.Red("\nError: %v", err)
		os.Exit(1)
	}
	color.Green("✓")

	fmt.Printf("→ Config: %s\n", mgr.GetConfigPath())
	fmt.Println()

	return mgr, cfg
}

// getEditor returns the user's preferred editor
func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	return "nano"
}
