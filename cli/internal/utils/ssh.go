package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pressctl/cli/pkg/models"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// TestSSHConnection tests SSH connectivity to a server
func TestSSHConnection(server models.Server) error {
	// Expand home directory in key file path
	keyFile := server.SSH.KeyFile
	if strings.HasPrefix(keyFile, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		keyFile = filepath.Join(homeDir, keyFile[1:])
	}

	// Read SSH private key
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("failed to read SSH key file %s: %w", keyFile, err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH private key: %w", err)
	}

	// Configure SSH client with TOFU host key verification
	// This validates against known_hosts if the file exists and the host is known,
	// or automatically accepts and saves unknown host keys
	config := &ssh.ClientConfig{
		User: server.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: trustOnFirstUseCallback(),
		Timeout:         10 * time.Second,
	}

	// Connect to server
	addr := fmt.Sprintf("%s:%d", server.IP, server.SSH.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH connection failed to %s: %w", addr, err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Test command execution
	output, err := session.CombinedOutput("echo 'pressctl-test'")
	if err != nil {
		return fmt.Errorf("test command failed: %w", err)
	}

	if strings.TrimSpace(string(output)) != "pressctl-test" {
		return fmt.Errorf("unexpected test output: %s", output)
	}

	return nil
}

// RunSSHCommand runs a command on a remote server over SSH and returns its combined output.
func RunSSHCommand(server models.Server, command string) (string, error) {
	keyFile := server.SSH.KeyFile
	if strings.HasPrefix(keyFile, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand home directory: %w", err)
		}
		keyFile = filepath.Join(homeDir, keyFile[1:])
	}

	key, err := os.ReadFile(keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read SSH key file %s: %w", keyFile, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to parse SSH private key: %w", err)
	}

	cfg := &ssh.ClientConfig{
		User: server.SSH.User,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: trustOnFirstUseCallback(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", server.IP, server.SSH.Port)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return "", fmt.Errorf("SSH connection failed: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(command)
	return string(out), err
}

// PortConflict describes a port that is already in use on the target server.
type PortConflict struct {
	Port    int
	Service string // e.g. "apache2", "mysql"
}

func (c PortConflict) String() string {
	return fmt.Sprintf("port %d is already in use by %s", c.Port, c.Service)
}

// CheckPortConflicts checks whether ports required by the provisioning playbook
// (80, 443, 3306) are already occupied on the remote server.
func CheckPortConflicts(server models.Server) ([]PortConflict, error) {
	// Use ss to list listening TCP ports; fall back to netstat if unavailable.
	cmd := `ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null || true`
	out, err := RunSSHCommand(server, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to check ports: %w", err)
	}

	portServices := map[int]string{
		80:   "web server (Apache/Caddy/other)",
		443:  "web server (Apache/Caddy/other)",
		3306: "MySQL/MariaDB",
	}

	var conflicts []PortConflict
	for port, service := range portServices {
		pattern := fmt.Sprintf(":%d ", port)
		if strings.Contains(out, pattern) {
			conflicts = append(conflicts, PortConflict{Port: port, Service: service})
		}
	}

	return conflicts, nil
}

// getHostKeyCallback returns a host key callback using the user's known_hosts file
func getHostKeyCallback() (ssh.HostKeyCallback, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	return knownhosts.New(knownHostsPath)
}

// trustOnFirstUseCallback returns a callback that accepts any host key
// and adds it to known_hosts on first connection (TOFU model)
func trustOnFirstUseCallback() ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// If we can't get home dir, just accept the key
			return nil
		}

		knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")

		// Try to read existing known_hosts
		callback, err := knownhosts.New(knownHostsPath)
		if err == nil {
			// File exists, check against it
			err = callback(hostname, remote, key)
			if err == nil {
				return nil // Key is known and matches
			}
			// Check if it's a KeyError (unknown host or key mismatch)
			if keyErr, ok := err.(*knownhosts.KeyError); ok {
				// If Want is not empty, it means we expected different keys (mismatch)
				if len(keyErr.Want) > 0 {
					return fmt.Errorf("host key mismatch for %s - possible security issue", hostname)
				}
				// Want is empty, so host is unknown - fall through to add it
			}
		}

		// Key not in known_hosts, add it (TOFU)
		// Ensure .ssh directory exists
		sshDir := filepath.Join(homeDir, ".ssh")
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			return nil // Accept key even if we can't save it
		}

		// Append to known_hosts
		f, err := os.OpenFile(knownHostsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil // Accept key even if we can't save it
		}
		defer f.Close()

		// Format the known_hosts line
		line := knownhosts.Line([]string{hostname}, key)
		if _, err := f.WriteString(line + "\n"); err != nil {
			return nil // Accept key even if we can't save it
		}

		return nil
	}
}
