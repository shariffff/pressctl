package config

import (
	"fmt"
	"os"
)

// SchemaVersion is the config schema version this binary expects.
// Bump this (e.g. "1.1") whenever the Config struct changes in a way
// that requires existing configs to be updated.
const SchemaVersion = "1.0"

type migrationFn func(*Config) error

var migrations = []struct {
	from string
	to   string
	run  migrationFn
}{
	// Add future migrations here, for example:
	// {"1.0", "1.1", migrate_1_0_to_1_1},
}

// MigrateIfNeeded checks whether cfg needs migration and, if so, runs
// all applicable migrations in sequence. It writes a .bak file before
// touching disk and saves the updated config via mgr on success.
// Returns true if any migration was applied.
func MigrateIfNeeded(mgr *Manager, cfg *Config) (bool, error) {
	if cfg.Version == SchemaVersion {
		return false, nil
	}

	// Back up the config before touching it.
	backupPath := mgr.GetConfigPath() + ".bak"
	data, err := os.ReadFile(mgr.GetConfigPath())
	if err != nil {
		return false, fmt.Errorf("failed to read config for backup: %w", err)
	}
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return false, fmt.Errorf("failed to backup config: %w", err)
	}

	current := cfg.Version
	applied := false

	for _, m := range migrations {
		if current != m.from {
			continue
		}
		if err := m.run(cfg); err != nil {
			os.Rename(backupPath, mgr.GetConfigPath()) // restore on failure
			return false, fmt.Errorf("migration %s→%s failed: %w", m.from, m.to, err)
		}
		cfg.Version = m.to
		current = m.to
		applied = true
	}

	if !applied {
		// No migration matched — version mismatch but nothing to run.
		// This can happen when downgrading; just leave the config as-is.
		os.Remove(backupPath)
		return false, nil
	}

	if err := mgr.Save(cfg); err != nil {
		os.Rename(backupPath, mgr.GetConfigPath()) // restore on failure
		return false, fmt.Errorf("failed to save migrated config: %w", err)
	}

	os.Remove(backupPath)
	return true, nil
}
