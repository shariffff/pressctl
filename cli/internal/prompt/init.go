package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InitInput holds the input for init setup
type InitInput struct {
	SSHPublicKey string
}

// PromptInitSetup auto-detects SSH public key for initial setup
func PromptInitSetup() (*InitInput, error) {
	input := &InitInput{}

	// Auto-detect SSH public key
	sshPubKeys, err := findSSHPublicKeys()
	if err != nil || len(sshPubKeys) == 0 {
		return nil, fmt.Errorf("no SSH public keys found in ~/.ssh/. Please generate one with: ssh-keygen")
	}

	// Use the first key found
	input.SSHPublicKey = sshPubKeys[0]
	fmt.Printf("→ Using SSH public key: %s\n", sshPubKeys[0])

	return input, nil
}

// findSSHPublicKeys looks for public SSH keys in ~/.ssh/
func findSSHPublicKeys() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Only include .pub files
		if !strings.HasSuffix(name, ".pub") {
			continue
		}

		keyPath := filepath.Join(sshDir, name)
		keys = append(keys, keyPath)
	}

	return keys, nil
}
