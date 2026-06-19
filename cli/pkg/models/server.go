package models

import "time"

// SSHConfig holds SSH connection details for a server
type SSHConfig struct {
	User string `yaml:"user" validate:"required"`
	Port int    `yaml:"port" validate:"required,min=1,max=65535"`
}

// ServerCredentials holds server-specific credentials
type ServerCredentials struct {
	MySQLWordmonbotPassword string `yaml:"mysql_pressctlbot_password,omitempty"`
}

// Server represents a managed server
type Server struct {
	Name          string            `yaml:"name" validate:"required"`
	Hostname      string            `yaml:"hostname" validate:"required"`
	IP            string            `yaml:"ip" validate:"required,ip"`
	SSH           SSHConfig         `yaml:"ssh"`
	Credentials   ServerCredentials `yaml:"credentials,omitempty"`
	PHPVersion    string            `yaml:"php_version,omitempty"`
	Stack         string            `yaml:"stack,omitempty" validate:"omitempty,oneof=traditional frankenphp"`
	Status        string            `yaml:"status" validate:"oneof=provisioned unprovisioned error"`
	ProvisionedAt *time.Time        `yaml:"provisioned_at,omitempty"`
	Sites         []Site            `yaml:"sites,omitempty"`
}

// SupportedPHPVersions lists PHP versions available for provisioning
var SupportedPHPVersions = []string{"8.5", "8.4", "8.3", "8.2", "8.1"}

// DefaultPHPVersion is the default PHP version for new servers
const DefaultPHPVersion = "8.3"

// Server stack identifiers. The stack determines how the web/runtime layer is
// provisioned and how sites are served.
const (
	// StackTraditional is host-installed Nginx + PHP-FPM (the original stack).
	StackTraditional = "traditional"
	// StackFrankenPHP is a Dockerized FrankenPHP container (Caddy + PHP) that
	// serves all sites and terminates TLS automatically.
	StackFrankenPHP = "frankenphp"
)

// DefaultStack is the stack used when a server does not specify one.
const DefaultStack = StackTraditional

// SupportedStacks lists the stacks available for provisioning.
var SupportedStacks = []string{StackTraditional, StackFrankenPHP}

// EffectiveStack returns the server's stack, falling back to DefaultStack when
// unset. This keeps configs written before stack selection behaving as
// traditional with no migration required.
func (s Server) EffectiveStack() string {
	if s.Stack == "" {
		return DefaultStack
	}
	return s.Stack
}

// IsValidStack reports whether the given stack identifier is supported.
func IsValidStack(stack string) bool {
	for _, v := range SupportedStacks {
		if v == stack {
			return true
		}
	}
	return false
}
