package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migrator handles database migrations
type Migrator struct {
	db             *sql.DB
	migrationsPath string
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Name    string
	Applied bool
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// CreateMigrationsTable creates the migrations tracking table
func (m *Migrator) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS public.schema_migrations (
			id SERIAL PRIMARY KEY,
			migration_name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// Up runs all pending migrations
func (m *Migrator) Up() ([]string, error) {
	// Get list of all migrations
	allMigrations, err := m.listMigrations()
	if err != nil {
		return nil, err
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	// Find pending migrations
	pendingMigrations := m.findPendingMigrations(allMigrations, appliedMigrations)

	if len(pendingMigrations) == 0 {
		return []string{}, nil
	}

	// Apply each pending migration
	applied := []string{}
	for _, migration := range pendingMigrations {
		if err := m.applyMigration(migration); err != nil {
			return applied, fmt.Errorf("failed to apply migration %s: %w", migration, err)
		}
		applied = append(applied, migration)
	}

	return applied, nil
}

// Down rolls back the last applied migration
func (m *Migrator) Down() (string, error) {
	// Get the last applied migration
	var migrationName string
	query := `
		SELECT migration_name
		FROM public.schema_migrations
		ORDER BY applied_at DESC
		LIMIT 1
	`

	err := m.db.QueryRow(query).Scan(&migrationName)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get last migration: %w", err)
	}

	// Rollback the migration
	if err := m.rollbackMigration(migrationName); err != nil {
		return "", fmt.Errorf("failed to rollback migration %s: %w", migrationName, err)
	}

	return migrationName, nil
}

// Status returns the status of all migrations
func (m *Migrator) Status() ([]MigrationStatus, error) {
	// Get all migrations
	allMigrations, err := m.listMigrations()
	if err != nil {
		return nil, err
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	// Build status list
	appliedMap := make(map[string]bool)
	for _, name := range appliedMigrations {
		appliedMap[name] = true
	}

	status := make([]MigrationStatus, len(allMigrations))
	for i, name := range allMigrations {
		status[i] = MigrationStatus{
			Name:    name,
			Applied: appliedMap[name],
		}
	}

	return status, nil
}

// Create creates a new migration with the given name
func (m *Migrator) Create(name string) (string, error) {
	// Generate migration name with timestamp
	timestamp := time.Now().Format("2006-01-02")
	migrationName := fmt.Sprintf("%s-%s", timestamp, strings.ToLower(strings.ReplaceAll(name, " ", "-")))

	// Create migration directory
	migrationDir := filepath.Join(m.migrationsPath, migrationName)
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create migration directory: %w", err)
	}

	// Create up.sql
	upFile := filepath.Join(migrationDir, "up.sql")
	upContent := fmt.Sprintf("-- Migration: %s (UP)\n-- Description: %s\n-- Date: %s\n\n-- Add your migration SQL here\n",
		migrationName, name, timestamp)
	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create up.sql: %w", err)
	}

	// Create down.sql
	downFile := filepath.Join(migrationDir, "down.sql")
	downContent := fmt.Sprintf("-- Migration: %s (DOWN)\n-- Description: Rollback %s\n-- Date: %s\n\n-- Add your rollback SQL here\n",
		migrationName, name, timestamp)
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create down.sql: %w", err)
	}

	return migrationName, nil
}

// listMigrations lists all migration directories
func (m *Migrator) listMigrations() ([]string, error) {
	entries, err := os.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Check if up.sql exists
			upFile := filepath.Join(m.migrationsPath, entry.Name(), "up.sql")
			if _, err := os.Stat(upFile); err == nil {
				migrations = append(migrations, entry.Name())
			}
		}
	}

	// Sort migrations alphabetically (which sorts by date due to naming convention)
	sort.Strings(migrations)

	return migrations, nil
}

// getAppliedMigrations gets list of applied migrations from database
func (m *Migrator) getAppliedMigrations() ([]string, error) {
	query := `
		SELECT migration_name
		FROM public.schema_migrations
		ORDER BY applied_at ASC
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan migration name: %w", err)
		}
		migrations = append(migrations, name)
	}

	return migrations, nil
}

// findPendingMigrations finds migrations that haven't been applied
func (m *Migrator) findPendingMigrations(all []string, applied []string) []string {
	appliedMap := make(map[string]bool)
	for _, name := range applied {
		appliedMap[name] = true
	}

	var pending []string
	for _, name := range all {
		if !appliedMap[name] {
			pending = append(pending, name)
		}
	}

	return pending
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(name string) error {
	// Read migration file
	upFile := filepath.Join(m.migrationsPath, name, "up.sql")
	content, err := os.ReadFile(upFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration
	query := `INSERT INTO public.schema_migrations (migration_name) VALUES ($1)`
	if _, err := tx.Exec(query, name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// rollbackMigration rolls back a single migration
func (m *Migrator) rollbackMigration(name string) error {
	// Read rollback file
	downFile := filepath.Join(m.migrationsPath, name, "down.sql")
	content, err := os.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("failed to read rollback file: %w", err)
	}

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute rollback
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute rollback SQL: %w", err)
	}

	// Remove migration record
	query := `DELETE FROM public.schema_migrations WHERE migration_name = $1`
	if _, err := tx.Exec(query, name); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
