package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pressctl/cli/internal/utils"
	"github.com/pressctl/cli/pkg/models"
)

// ServerInput holds the input for server creation
type ServerInput struct {
	Name       string
	Hostname   string
	IP         string
	SSHUser    string
	SSHPort    int
	SSHKey     string
	PHPVersion string
}

// PromptServerAdd prompts for server details
func PromptServerAdd() (*ServerInput, error) {
	input := &ServerInput{}

	// SSH key selection (conditional — must happen before the main form)
	sshKeys, err := findSSHKeys()
	if err != nil || len(sshKeys) == 0 {
		homeDir, _ := os.UserHomeDir()
		input.SSHKey = filepath.Join(homeDir, ".ssh", "id_rsa")
	} else if len(sshKeys) == 1 {
		input.SSHKey = sshKeys[0]
		fmt.Printf("Using SSH key: %s\n", sshKeys[0])
	} else {
		opts := make([]huh.Option[string], len(sshKeys))
		for i, k := range sshKeys {
			opts[i] = huh.NewOption(k, k)
		}
		if err := huh.NewSelect[string]().
			Title("SSH private key").
			Description("Select the key to use for authentication").
			Options(opts...).
			Value(&input.SSHKey).
			Run(); err != nil {
			return nil, normalizeErr(err)
		}
	}

	// Main form: name, IP, SSH user, SSH port, PHP version
	var portStr string
	phpVersions := make([]huh.Option[string], len(models.SupportedPHPVersions))
	for i, v := range models.SupportedPHPVersions {
		phpVersions[i] = huh.NewOption(v, v)
	}

	if err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server name").
				Description("A friendly name to identify this server (e.g., production-1)").
				Value(&input.Name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("server name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("IP address").
				Description("The server's public IP from your cloud provider").
				Value(&input.IP).
				Validate(utils.ValidateIP),
			huh.NewInput().
				Title("SSH user").
				Description("User for initial provisioning").
				Value(&input.SSHUser).
				Placeholder("root"),
			huh.NewInput().
				Title("SSH port").
				Value(&portStr).
				Placeholder("22").
				Validate(func(s string) error {
					if s == "" {
						return nil // allow empty, defaults to 22
					}
					return utils.ValidatePort(s)
				}),
			huh.NewSelect[string]().
				Title("PHP version").
				Options(phpVersions...).
				Value(&input.PHPVersion),
		),
	).Run(); err != nil {
		return nil, normalizeErr(err)
	}

	// Apply defaults
	if strings.TrimSpace(input.SSHUser) == "" {
		input.SSHUser = "root"
	}
	if portStr == "" {
		input.SSHPort = 22
	} else {
		input.SSHPort, _ = strconv.Atoi(portStr)
	}
	if input.PHPVersion == "" {
		input.PHPVersion = models.DefaultPHPVersion
	}
	input.Hostname = input.IP

	// Confirmation
	if err := confirmServerAdd(input); err != nil {
		return nil, err
	}

	return input, nil
}

// findSSHKeys looks for private SSH keys in ~/.ssh/
func findSSHKeys() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(filepath.Join(homeDir, ".ssh"))
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".pub") ||
			name == "known_hosts" || name == "known_hosts.old" ||
			name == "config" || name == "authorized_keys" || name == "environment" {
			continue
		}
		keyPath := filepath.Join(homeDir, ".ssh", name)
		content, err := os.ReadFile(keyPath)
		if err != nil {
			continue
		}
		if strings.HasPrefix(string(content), "-----BEGIN") {
			keys = append(keys, keyPath)
		}
	}

	return keys, nil
}

// ToServer converts ServerInput to models.Server
func (si *ServerInput) ToServer() models.Server {
	phpVersion := si.PHPVersion
	if phpVersion == "" {
		phpVersion = models.DefaultPHPVersion
	}
	return models.Server{
		Name:     si.Name,
		Hostname: si.Hostname,
		IP:       si.IP,
		SSH: models.SSHConfig{
			User:    si.SSHUser,
			Port:    si.SSHPort,
			KeyFile: si.SSHKey,
		},
		PHPVersion: phpVersion,
		Status:     "unprovisioned",
		Sites:      []models.Site{},
	}
}

func confirmServerAdd(input *ServerInput) error {
	fmt.Println()
	fmt.Printf("  Name:        %s\n", input.Name)
	fmt.Printf("  IP:          %s\n", input.IP)
	fmt.Printf("  SSH Key:     %s\n", input.SSHKey)
	fmt.Printf("  SSH User:    %s\n", input.SSHUser)
	fmt.Printf("  SSH Port:    %d\n", input.SSHPort)
	fmt.Printf("  PHP Version: %s\n", input.PHPVersion)
	fmt.Println()

	var confirm bool
	if err := huh.NewConfirm().
		Title("Provision this server?").
		Affirmative("Yes, provision").
		Negative("Cancel").
		Value(&confirm).
		Run(); err != nil {
		return normalizeErr(err)
	}
	if !confirm {
		return fmt.Errorf("cancelled")
	}
	return nil
}
