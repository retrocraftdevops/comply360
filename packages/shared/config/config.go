package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for a service
type Config struct {
	// Environment
	Environment string
	ServiceName string
	ServicePort string

	// Database
	DatabaseURL      string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseSSLMode  string
	DatabaseMaxConns int
	DatabaseMinConns int

	// Redis
	RedisURL      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// RabbitMQ
	RabbitMQURL      string
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQVHost    string

	// MinIO/S3
	MinIOEndpoint   string
	MinIOAccessKey  string
	MinIOSecretKey  string
	MinIOBucket     string
	MinIOUseSSL     bool

	// Odoo
	OdooURL      string
	OdooDatabase string
	OdooUsername string
	OdooPassword string

	// JWT
	JWTSecret              string
	JWTAccessTokenExpiry   time.Duration
	JWTRefreshTokenExpiry  time.Duration

	// Security
	MaxFailedLoginAttempts int
	AccountLockDuration    time.Duration
	BCryptCost             int

	// Rate Limiting
	RateLimitEnabled         bool
	RateLimitRequestsPerMin  int
	RateLimitBurstSize       int

	// CORS
	CORSAllowedOrigins   []string
	CORSAllowCredentials bool

	// Email/SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// SMS
	SMSProvider   string // twilio, africastalking
	SMSAPIKey     string
	SMSAPISecret  string
	SMSFromNumber string

	// Feature Flags
	EnableRegistration      bool
	EnableEmailVerification bool
	EnableMFA               bool

	// Monitoring
	MetricsEnabled bool
	MetricsPort    string
	TracingEnabled bool
	LogLevel       string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// Environment
		Environment: getEnv("APP_ENV", "development"),
		ServiceName: getEnv("SERVICE_NAME", "comply360-service"),
		ServicePort: getEnv("SERVICE_PORT", "8080"),

		// Database
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://comply360_user:comply360_password_dev@localhost:5432/comply360_db?sslmode=disable"),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseName:     getEnv("DATABASE_NAME", "comply360_db"),
		DatabaseUser:     getEnv("DATABASE_USER", "comply360_user"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "comply360_password_dev"),
		DatabaseSSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
		DatabaseMaxConns: getEnvInt("DATABASE_MAX_CONNS", 25),
		DatabaseMinConns: getEnvInt("DATABASE_MIN_CONNS", 5),

		// Redis
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		// RabbitMQ
		RabbitMQURL:      getEnv("RABBITMQ_URL", "amqp://comply360:dev_password@localhost:5672/"),
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnv("RABBITMQ_USER", "comply360"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "dev_password"),
		RabbitMQVHost:    getEnv("RABBITMQ_VHOST", "/"),

		// MinIO/S3
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "comply360"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "dev_password"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "comply360-documents"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),

		// Odoo
		OdooURL:      getEnv("ODOO_URL", "http://localhost:8069"),
		OdooDatabase: getEnv("ODOO_DATABASE", "comply360_odoo"),
		OdooUsername: getEnv("ODOO_USERNAME", "admin"),
		OdooPassword: getEnv("ODOO_PASSWORD", "admin"),

		// JWT
		JWTSecret:             getEnv("JWT_SECRET", "dev_secret_key_change_in_production"),
		JWTAccessTokenExpiry:  getEnvDuration("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		JWTRefreshTokenExpiry: getEnvDuration("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),

		// Security
		MaxFailedLoginAttempts: getEnvInt("MAX_FAILED_LOGIN_ATTEMPTS", 5),
		AccountLockDuration:    getEnvDuration("ACCOUNT_LOCK_DURATION", 30*time.Minute),
		BCryptCost:             getEnvInt("BCRYPT_COST", 12),

		// Rate Limiting
		RateLimitEnabled:        getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitRequestsPerMin: getEnvInt("RATE_LIMIT_REQUESTS_PER_MIN", 1000),
		RateLimitBurstSize:      getEnvInt("RATE_LIMIT_BURST_SIZE", 100),

		// CORS
		CORSAllowedOrigins:   getEnvSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),

		// Email/SMTP
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@comply360.com"),

		// SMS
		SMSProvider:   getEnv("SMS_PROVIDER", ""), // twilio, africastalking
		SMSAPIKey:     getEnv("SMS_API_KEY", ""),
		SMSAPISecret:  getEnv("SMS_API_SECRET", ""),
		SMSFromNumber: getEnv("SMS_FROM_NUMBER", ""),

		// Feature Flags
		EnableRegistration:      getEnvBool("ENABLE_REGISTRATION", true),
		EnableEmailVerification: getEnvBool("ENABLE_EMAIL_VERIFICATION", true),
		EnableMFA:               getEnvBool("ENABLE_MFA", false),

		// Monitoring
		MetricsEnabled: getEnvBool("METRICS_ENABLED", false),
		MetricsPort:    getEnv("METRICS_PORT", "9090"),
		TracingEnabled: getEnvBool("TRACING_ENABLED", false),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.IsProduction() {
		if c.JWTSecret == "dev_secret_key_change_in_production" {
			return fmt.Errorf("JWT_SECRET must be changed in production")
		}
		if c.DatabasePassword == "comply360_password_dev" {
			return fmt.Errorf("DATABASE_PASSWORD must be changed in production")
		}
	}
	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma
		var result []string
		for _, item := range splitAndTrim(value, ",") {
			if item != "" {
				result = append(result, item)
			}
		}
		return result
	}
	return defaultValue
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, item := range splitString(s, sep) {
		trimmed := trimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	// Simple split implementation
	var result []string
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, current)
			current = ""
			i += len(sep) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
