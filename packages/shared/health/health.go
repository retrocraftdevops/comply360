package health

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Status represents health check status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusDegraded  Status = "degraded"
	StatusUnhealthy Status = "unhealthy"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Latency   string                 `json:"latency,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// HealthCheck represents overall health of the service
type HealthCheck struct {
	Service    string         `json:"service"`
	Version    string         `json:"version"`
	Status     Status         `json:"status"`
	Checks     []CheckResult  `json:"checks"`
	Timestamp  time.Time      `json:"timestamp"`
}

// Checker performs health checks
type Checker struct {
	serviceName    string
	serviceVersion string
	checks         []func() CheckResult
}

// NewChecker creates a new health checker
func NewChecker(serviceName, serviceVersion string) *Checker {
	return &Checker{
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		checks:         []func() CheckResult{},
	}
}

// AddCheck adds a custom check function
func (c *Checker) AddCheck(checkFunc func() CheckResult) {
	c.checks = append(c.checks, checkFunc)
}

// AddDatabaseCheck adds a database health check
func (c *Checker) AddDatabaseCheck(db *sql.DB) {
	c.AddCheck(func() CheckResult {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		result := CheckResult{
			Name:      "database",
			Timestamp: time.Now(),
		}

		if err := db.PingContext(ctx); err != nil {
			result.Status = StatusUnhealthy
			result.Message = fmt.Sprintf("Database ping failed: %v", err)
			return result
		}

		// Check connection stats
		stats := db.Stats()
		result.Status = StatusHealthy
		result.Latency = time.Since(start).String()
		result.Metadata = map[string]interface{}{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		}

		return result
	})
}

// AddRedisCheck adds a Redis health check
func (c *Checker) AddRedisCheck(client *redis.Client) {
	c.AddCheck(func() CheckResult {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		result := CheckResult{
			Name:      "redis",
			Timestamp: time.Now(),
		}

		if err := client.Ping(ctx).Err(); err != nil {
			result.Status = StatusUnhealthy
			result.Message = fmt.Sprintf("Redis ping failed: %v", err)
			return result
		}

		result.Status = StatusHealthy
		result.Latency = time.Since(start).String()

		return result
	})
}

// AddRabbitMQCheck adds a RabbitMQ health check
func (c *Checker) AddRabbitMQCheck(conn *amqp.Connection) {
	c.AddCheck(func() CheckResult {
		result := CheckResult{
			Name:      "rabbitmq",
			Timestamp: time.Now(),
		}

		if conn == nil || conn.IsClosed() {
			result.Status = StatusUnhealthy
			result.Message = "RabbitMQ connection is closed"
			return result
		}

		result.Status = StatusHealthy
		result.Message = "Connected"

		return result
	})
}

// Check runs all health checks and returns the overall health
func (c *Checker) Check() HealthCheck {
	health := HealthCheck{
		Service:   c.serviceName,
		Version:   c.serviceVersion,
		Status:    StatusHealthy,
		Checks:    []CheckResult{},
		Timestamp: time.Now(),
	}

	// Run all checks
	for _, checkFunc := range c.checks {
		result := checkFunc()
		health.Checks = append(health.Checks, result)

		// Update overall status
		if result.Status == StatusUnhealthy {
			health.Status = StatusUnhealthy
		} else if result.Status == StatusDegraded && health.Status != StatusUnhealthy {
			health.Status = StatusDegraded
		}
	}

	// If no checks were run, service is healthy
	if len(health.Checks) == 0 {
		health.Status = StatusHealthy
	}

	return health
}

// IsHealthy returns true if the service is healthy
func (h *HealthCheck) IsHealthy() bool {
	return h.Status == StatusHealthy
}

// IsDegraded returns true if the service is degraded
func (h *HealthCheck) IsDegraded() bool {
	return h.Status == StatusDegraded
}

// IsUnhealthy returns true if the service is unhealthy
func (h *HealthCheck) IsUnhealthy() bool {
	return h.Status == StatusUnhealthy
}
