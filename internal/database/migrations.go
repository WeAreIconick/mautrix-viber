// Package database migrations provides database migration tool for schema updates.
package database

import (
	"fmt"
)

// Migration represents a database migration.
type Migration struct {
	Version int
	Up      string // SQL for upgrading
	Down    string // SQL for downgrading
}

// Migrator manages database migrations.
type Migrator struct {
	db        *DB
	migrations []Migration
}

// NewMigrator creates a new migrator.
func NewMigrator(db *DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []Migration{},
	}
}

// RegisterMigration registers a migration.
func (m *Migrator) RegisterMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// GetCurrentVersion gets the current database version.
func (m *Migrator) GetCurrentVersion() (int, error) {
	// Check if migrations table exists
	var version int
	err := m.db.db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err != nil {
		// No migrations table - assume version 0
		return 0, nil
	}
	return version, nil
}

// Migrate runs migrations up to the target version.
func (m *Migrator) Migrate(targetVersion int) error {
	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("get current version: %w", err)
	}
	
	// Create migrations table if it doesn't exist
	_, err = m.db.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}
	
	// Apply migrations
	for _, migration := range m.migrations {
		if migration.Version > currentVersion && migration.Version <= targetVersion {
			if err := m.applyMigration(migration); err != nil {
				return fmt.Errorf("apply migration %d: %w", migration.Version, err)
			}
		}
	}
	
	return nil
}

// applyMigration applies a single migration.
func (m *Migrator) applyMigration(migration Migration) error {
	// Run migration SQL
	if _, err := m.db.db.Exec(migration.Up); err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}
	
	// Record migration
	_, err := m.db.db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration.Version)
	return err
}

// Rollback rolls back migrations down to the target version.
func (m *Migrator) Rollback(targetVersion int) error {
	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("get current version: %w", err)
	}
	
	// Rollback migrations in reverse order
	for i := len(m.migrations) - 1; i >= 0; i-- {
		migration := m.migrations[i]
		if migration.Version <= currentVersion && migration.Version > targetVersion {
			if err := m.rollbackMigration(migration); err != nil {
				return fmt.Errorf("rollback migration %d: %w", migration.Version, err)
			}
		}
	}
	
	return nil
}

// rollbackMigration rolls back a single migration.
func (m *Migrator) rollbackMigration(migration Migration) error {
	if migration.Down == "" {
		return fmt.Errorf("migration %d has no down migration", migration.Version)
	}
	
	// Run rollback SQL
	if _, err := m.db.db.Exec(migration.Down); err != nil {
		return fmt.Errorf("execute rollback: %w", err)
	}
	
	// Remove migration record
	_, err := m.db.db.Exec("DELETE FROM schema_migrations WHERE version = ?", migration.Version)
	return err
}

