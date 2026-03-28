package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

// Validator handles configuration validation
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new config validator
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateStruct performs struct-level validation using tags
func (v *Validator) ValidateStruct(config *Config) error {
	if err := v.validate.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

// ValidateBusinessRules performs business logic validation
func (v *Validator) ValidateBusinessRules(config *Config) error {
	// Check unique server names
	serverNames := make(map[string]bool)
	for _, server := range config.Servers {
		if serverNames[server.Name] {
			return fmt.Errorf("duplicate server name: %s", server.Name)
		}
		serverNames[server.Name] = true
	}

	// Check unique domains across all servers
	domains := make(map[string]string)
	for _, server := range config.Servers {
		for _, site := range server.Sites {
			for _, domain := range site.Domains {
				if existingServer, exists := domains[domain.Domain]; exists {
					return fmt.Errorf("domain %s exists on both server %s and %s",
						domain.Domain, existingServer, server.Name)
				}
				domains[domain.Domain] = server.Name
			}
		}
	}

	return nil
}

// ValidateAnsibleEnvironment checks if Ansible and required files exist
func (v *Validator) ValidateAnsibleEnvironment(config *Config) error {
	// Check if ansible-playbook exists in PATH
	if _, err := exec.LookPath("ansible-playbook"); err != nil {
		return fmt.Errorf("ansible-playbook not found in PATH. Please install Ansible")
	}

	// Expand home directory if path starts with ~
	ansiblePath := config.Ansible.Path
	if len(ansiblePath) > 0 && ansiblePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		ansiblePath = filepath.Join(homeDir, ansiblePath[1:])
	}

	// Check if ansible path exists
	if _, err := os.Stat(ansiblePath); os.IsNotExist(err) {
		return fmt.Errorf("ansible path does not exist: %s", ansiblePath)
	}

	// Check if required playbooks exist
	requiredPlaybooks := []string{
		"provision.yml",
		"website.yml",
	}

	for _, playbook := range requiredPlaybooks {
		path := filepath.Join(ansiblePath, playbook)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("required playbook not found: %s", path)
		}
	}

	// Check if roles directory exists
	rolesPath := filepath.Join(ansiblePath, "roles")
	if _, err := os.Stat(rolesPath); os.IsNotExist(err) {
		return fmt.Errorf("roles directory not found: %s", rolesPath)
	}

	// Check if group_vars/all.yml exists
	allVarsPath := filepath.Join(ansiblePath, "group_vars", "all.yml")
	if _, err := os.Stat(allVarsPath); os.IsNotExist(err) {
		return fmt.Errorf("group_vars/all.yml not found: %s", allVarsPath)
	}

	return nil
}

// Validate runs all validation checks
func (v *Validator) Validate(config *Config) error {
	if err := v.ValidateStruct(config); err != nil {
		return err
	}

	if err := v.ValidateBusinessRules(config); err != nil {
		return err
	}

	if err := v.ValidateAnsibleEnvironment(config); err != nil {
		return err
	}

	return nil
}
