package prompt

import (
	"fmt"
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
	PHPVersion string
}

// PromptServerAdd prompts for server details
func PromptServerAdd() (*ServerInput, error) {
	input := &ServerInput{}

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
			User: si.SSHUser,
			Port: si.SSHPort,
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
