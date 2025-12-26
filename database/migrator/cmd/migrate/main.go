package main

import (
	"fmt"
	"os"

	"github.com/comply360/migrator/internal/db"
	"github.com/comply360/migrator/internal/migrations"
	"github.com/spf13/cobra"
)

var (
	dbURL          string
	migrationsPath string
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration tool for Comply360",
	Long:  `A CLI tool to manage database migrations for the Comply360 multi-tenant platform.`,
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	RunE:  runUp,
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	RunE:  runDown,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE:  runStatus,
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreate,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&dbURL, "db-url", getEnv("DATABASE_URL", ""), "Database connection URL")
	rootCmd.PersistentFlags().StringVar(&migrationsPath, "path", "./migrations", "Path to migrations directory")

	// Add commands
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(createCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runUp(cmd *cobra.Command, args []string) error {
	database, err := db.Connect(dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	migrator := migrations.NewMigrator(database, migrationsPath)

	if err := migrator.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	applied, err := migrator.Up()
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if len(applied) == 0 {
		fmt.Println("No pending migrations")
	} else {
		fmt.Printf("Successfully applied %d migration(s):\n", len(applied))
		for _, name := range applied {
			fmt.Printf("  ✓ %s\n", name)
		}
	}

	return nil
}

func runDown(cmd *cobra.Command, args []string) error {
	database, err := db.Connect(dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	migrator := migrations.NewMigrator(database, migrationsPath)

	name, err := migrator.Down()
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	if name == "" {
		fmt.Println("No migrations to rollback")
	} else {
		fmt.Printf("Successfully rolled back: %s\n", name)
	}

	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	database, err := db.Connect(dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	migrator := migrations.NewMigrator(database, migrationsPath)

	if err := migrator.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	status, err := migrator.Status()
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")
	for _, s := range status {
		statusIcon := "✗"
		if s.Applied {
			statusIcon = "✓"
		}
		fmt.Printf("%s %s\n", statusIcon, s.Name)
	}

	return nil
}

func runCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	migrator := migrations.NewMigrator(nil, migrationsPath)

	created, err := migrator.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	fmt.Printf("Created migration: %s\n", created)
	fmt.Printf("  - %s/up.sql\n", created)
	fmt.Printf("  - %s/down.sql\n", created)

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
