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
	Status        string            `yaml:"status" validate:"oneof=provisioned unprovisioned error"`
	ProvisionedAt *time.Time        `yaml:"provisioned_at,omitempty"`
	Sites         []Site            `yaml:"sites,omitempty"`
}

// SupportedPHPVersions lists PHP versions available for provisioning
var SupportedPHPVersions = []string{"8.5", "8.4", "8.3", "8.2", "8.1"}

// DefaultPHPVersion is the default PHP version for new servers
const DefaultPHPVersion = "8.3"
