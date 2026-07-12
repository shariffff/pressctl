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
	Name              string
	Hostname          string
	IP                string
	SSHUser           string
	SSHPort           int
	PHPVersion        string
	Stack             string
	FrankenPHPRuntime string
}

// PromptServerAdd prompts for server details
func PromptServerAdd() (*ServerInput, error) {
	input := &ServerInput{}

	// Main form: name, IP, SSH user, SSH port, stack, PHP version
	var portStr string
	phpVersions := make([]huh.Option[string], len(models.SupportedPHPVersions))
	for i, v := range models.SupportedPHPVersions {
		phpVersions[i] = huh.NewOption(v, v)
	}

	// The stack select offers three choices. FrankenPHP has two runtimes, so we
	// encode the choice as "stack" or "stack:runtime" and split it afterwards.
	const (
		choiceTraditional   = models.StackTraditional
		choiceFrankenNative = models.StackFrankenPHP + ":" + models.RuntimeNative
		choiceFrankenDocker = models.StackFrankenPHP + ":" + models.RuntimeDocker
	)
	var stackChoice string
	stacks := []huh.Option[string]{
		huh.NewOption("Traditional (Nginx + PHP-FPM)", choiceTraditional),
		huh.NewOption("FrankenPHP — native binary (Caddy auto-HTTPS, no Docker)", choiceFrankenNative),
		huh.NewOption("FrankenPHP — Docker container (Caddy auto-HTTPS)", choiceFrankenDocker),
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
				Title("Server stack").
				Description("Traditional installs Nginx + PHP-FPM on the host. FrankenPHP (Caddy + PHP, automatic HTTPS) runs either as a native systemd binary or in Docker").
				Options(stacks...).
				Value(&stackChoice),
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

	// Split the composite stack choice into stack + FrankenPHP runtime.
	if before, after, found := strings.Cut(stackChoice, ":"); found {
		input.Stack = before
		input.FrankenPHPRuntime = after
	} else {
		input.Stack = stackChoice
	}
	if input.Stack == "" {
		input.Stack = models.DefaultStack
	}
	if input.Stack == models.StackFrankenPHP && input.FrankenPHPRuntime == "" {
		input.FrankenPHPRuntime = models.DefaultFrankenPHPRuntime
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
	stack := si.Stack
	if stack == "" {
		stack = models.DefaultStack
	}
	runtime := si.FrankenPHPRuntime
	if stack == models.StackFrankenPHP && runtime == "" {
		runtime = models.DefaultFrankenPHPRuntime
	}
	if stack != models.StackFrankenPHP {
		runtime = "" // runtime only applies to the frankenphp stack
	}
	return models.Server{
		Name:     si.Name,
		Hostname: si.Hostname,
		IP:       si.IP,
		SSH: models.SSHConfig{
			User: si.SSHUser,
			Port: si.SSHPort,
		},
		PHPVersion:        phpVersion,
		Stack:             stack,
		FrankenPHPRuntime: runtime,
		Status:            "unprovisioned",
		Sites:             []models.Site{},
	}
}

func confirmServerAdd(input *ServerInput) error {
	fmt.Println()
	fmt.Printf("  Name:        %s\n", input.Name)
	fmt.Printf("  IP:          %s\n", input.IP)
	fmt.Printf("  SSH User:    %s\n", input.SSHUser)
	fmt.Printf("  SSH Port:    %d\n", input.SSHPort)
	if input.Stack == models.StackFrankenPHP {
		fmt.Printf("  Stack:       %s (%s runtime)\n", input.Stack, input.FrankenPHPRuntime)
	} else {
		fmt.Printf("  Stack:       %s\n", input.Stack)
	}
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
