package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pressctl/cli/pkg/models"
)

// sshArgs builds the common ssh command arguments for a server.
func sshArgs(server models.Server) []string {
	return []string{
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "BatchMode=yes",
		"-o", "ConnectTimeout=10",
		"-p", fmt.Sprintf("%d", server.SSH.Port),
		fmt.Sprintf("%s@%s", server.SSH.User, server.IP),
	}
}

// TestSSHConnection tests SSH connectivity to a server using the system ssh command.
func TestSSHConnection(server models.Server) error {
	args := append(sshArgs(server), "echo", "pressctl-test")
	cmd := exec.Command("ssh", args...)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("SSH connection failed to %s@%s:%d: %v",
			server.SSH.User, server.IP, server.SSH.Port, err)
	}

	if strings.TrimSpace(string(output)) != "pressctl-test" {
		return fmt.Errorf("unexpected test output: %s", output)
	}

	return nil
}

// RunSSHCommand runs a command on a remote server over SSH and returns its stdout.
func RunSSHCommand(server models.Server, command string) (string, error) {
	args := append(sshArgs(server), command)
	cmd := exec.Command("ssh", args...)
	out, err := cmd.Output()
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
	out, err := RunSSHCommand(server, `ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null || true`)
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
