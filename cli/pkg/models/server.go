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
	PHPVersions   []string          `yaml:"php_versions,omitempty"`
	Status        string            `yaml:"status" validate:"oneof=provisioned unprovisioned error"`
	ProvisionedAt *time.Time        `yaml:"provisioned_at,omitempty"`
	Sites         []Site            `yaml:"sites,omitempty"`
}

// SupportedPHPVersions lists PHP versions available for provisioning
var SupportedPHPVersions = []string{"8.5", "8.4", "8.3", "8.2", "8.1"}

// DefaultPHPVersion is the default PHP version for new servers
const DefaultPHPVersion = "8.3"

// IsValidPHPVersion reports whether the given PHP version is supported.
func IsValidPHPVersion(version string) bool {
	for _, v := range SupportedPHPVersions {
		if v == version {
			return true
		}
	}
	return false
}

// HasPHPVersion reports whether the given PHP version is installed on the
// server. The server's default PHPVersion is always considered installed.
func (s *Server) HasPHPVersion(version string) bool {
	if s.PHPVersion == version {
		return true
	}
	for _, v := range s.PHPVersions {
		if v == version {
			return true
		}
	}
	return false
}

// InstalledPHPVersions returns the list of PHP versions installed on the
// server, always including the server's default PHPVersion.
func (s *Server) InstalledPHPVersions() []string {
	result := []string{}
	if s.PHPVersion != "" {
		result = append(result, s.PHPVersion)
	}
	for _, v := range s.PHPVersions {
		if v != s.PHPVersion {
			result = append(result, v)
		}
	}
	return result
}
