package testing

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// TestDB holds test database connection
type TestDB struct {
	DB       *sql.DB
	Schema   string
	TenantID uuid.UUID
}

// getTestDBURL returns the test database URL from environment or uses defaults
func getTestDBURL() string {
	host := getEnv("APP_DB_HOST", "localhost")
	port := getEnv("APP_DB_PORT", "5432")
	user := getEnv("APP_DB_USER", "comply360_app_user")
	password := getEnv("APP_DB_PASSWORD", "comply360_app_secure_pass")
	dbname := getEnv("APP_DB_NAME", "comply360_app")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *TestDB {
	dbURL := getTestDBURL()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Create a test tenant schema
	tenantID := uuid.New()
	schema := fmt.Sprintf("test_%s", uuid.New().String()[:8])

	// Create schema
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	if err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}

	// Set search path
	_, err = db.Exec(fmt.Sprintf("SET search_path TO %s, public", schema))
	if err != nil {
		t.Fatalf("Failed to set search path: %v", err)
	}

	return &TestDB{
		DB:       db,
		Schema:   schema,
		TenantID: tenantID,
	}
}

// Cleanup drops the test schema and closes the connection
func (tdb *TestDB) Cleanup(t *testing.T) {
	if tdb.DB != nil {
		_, err := tdb.DB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", tdb.Schema))
		if err != nil {
			t.Logf("Warning: Failed to drop test schema: %v", err)
		}
		tdb.DB.Close()
	}
}

// CreateTestTables creates necessary tables in the test schema
func (tdb *TestDB) CreateTestTables(t *testing.T) {
	// Create users table
	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			email VARCHAR(255) NOT NULL,
			password_hash VARCHAR(255),
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			phone VARCHAR(50),
			mobile VARCHAR(50),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			email_verified BOOLEAN NOT NULL DEFAULT false,
			email_verified_at TIMESTAMP,
			mfa_enabled BOOLEAN NOT NULL DEFAULT false,
			mfa_method VARCHAR(50),
			mfa_secret VARCHAR(255),
			failed_login_attempts INT NOT NULL DEFAULT 0,
			locked_until TIMESTAMP,
			last_login_at TIMESTAMP,
			last_login_ip VARCHAR(50),
			password_changed_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Create user_roles table
	_, err = tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS user_roles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL,
			granted_by UUID,
			granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMP,
			UNIQUE(user_id, role)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create user_roles table: %v", err)
	}
}

// TestRedis holds test Redis connection
type TestRedis struct {
	Client *redis.Client
}

// SetupTestRedis creates a test Redis connection
func SetupTestRedis(t *testing.T) *TestRedis {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use DB 15 for testing
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
	}

	// Flush test database
	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("Failed to flush Redis test database: %v", err)
	}

	return &TestRedis{Client: client}
}

// Cleanup closes Redis connection and flushes test data
func (tr *TestRedis) Cleanup(t *testing.T) {
	if tr.Client != nil {
		ctx := context.Background()
		tr.Client.FlushDB(ctx)
		tr.Client.Close()
	}
}

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: %v", msgAndArgs[0], err)
		} else {
			t.Fatalf("Unexpected error: %v", err)
		}
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected error but got nil", msgAndArgs[0])
		} else {
			t.Fatal("Expected error but got nil")
		}
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected %v but got %v", msgAndArgs[0], expected, actual)
		} else {
			t.Fatalf("Expected %v but got %v", expected, actual)
		}
	}
}

// AssertNotEqual fails the test if expected == actual
func AssertNotEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if expected == actual {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected values to be different but both are %v", msgAndArgs[0], expected)
		} else {
			t.Fatalf("Expected values to be different but both are %v", expected)
		}
	}
}

// AssertNil fails the test if value is not nil
func AssertNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	// Handle typed nil pointers properly
	if value != nil && !isNil(value) {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected nil but got %v", msgAndArgs[0], value)
		} else {
			t.Fatalf("Expected nil but got %v", value)
		}
	}
}

// isNil checks if a value is nil using reflection to handle typed nils
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	// Use reflection to check if the value is a typed nil
	v := reflect.ValueOf(value)
	kind := v.Kind()
	return (kind == reflect.Ptr || kind == reflect.Interface || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan || kind == reflect.Func) && v.IsNil()
}

// AssertNotNil fails the test if value is nil
func AssertNotNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value == nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected non-nil value", msgAndArgs[0])
		} else {
			t.Fatal("Expected non-nil value but got nil")
		}
	}
}

// AssertTrue fails the test if condition is false
func AssertTrue(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if !condition {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected true but got false", msgAndArgs[0])
		} else {
			t.Fatal("Expected true but got false")
		}
	}
}

// AssertFalse fails the test if condition is true
func AssertFalse(t *testing.T, condition bool, msgAndArgs ...interface{}) {
	t.Helper()
	if condition {
		if len(msgAndArgs) > 0 {
			t.Fatalf("%s: expected false but got true", msgAndArgs[0])
		} else {
			t.Fatal("Expected false but got true")
		}
	}
}
