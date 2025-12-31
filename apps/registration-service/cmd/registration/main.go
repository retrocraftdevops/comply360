package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/comply360/registration-service/internal/events"
	"github.com/comply360/registration-service/internal/handlers"
	"github.com/comply360/registration-service/internal/repository"
	"github.com/comply360/registration-service/internal/services"
	sharedmiddleware "github.com/comply360/shared/middleware"
	sharedsentry "github.com/comply360/shared/sentry"
	"github.com/comply360/shared/websocket"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	// SECURITY: Default to SSL required for production-grade security
	dbURL := getEnv("DATABASE_URL", "postgresql://comply360:dev_password@localhost:5432/comply360_dev?sslmode=require")
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://comply360:dev_password@localhost:5672/")
	port := getEnv("REGISTRATION_SERVICE_PORT", "8083")
	jwtSecret := os.Getenv("JWT_SECRET")

	// SECURITY: Enforce JWT_SECRET configuration
	if jwtSecret == "" {
		log.Fatal("FATAL: JWT_SECRET environment variable is not set. Application cannot start without a secure JWT secret.")
	}
	if jwtSecret == "changeme" || jwtSecret == "changeme-change-this-in-production" || jwtSecret == "dev_secret_key_change_in_production" {
		log.Fatal("FATAL: JWT_SECRET is set to a default/insecure value. Please configure a secure secret.")
	}
	if len(jwtSecret) < 32 {
		log.Fatalf("FATAL: JWT_SECRET must be at least 32 characters long (current length: %d)", len(jwtSecret))
	}
	log.Println("JWT_SECRET validated successfully")

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool for production performance
	db.SetMaxOpenConns(25)                      // Maximum number of open connections to the database
	db.SetMaxIdleConns(5)                       // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute)      // Connection lifetime: 5 minutes
	db.SetConnMaxIdleTime(10 * time.Minute)     // Idle connection timeout: 10 minutes

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL with optimized connection pooling")

	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	// PRODUCTION: Initialize Sentry for error tracking and monitoring
	env := getEnv("APP_ENV", "development")
	sentryDSN := os.Getenv("SENTRY_DSN")
	if err := sharedsentry.InitSentry(sharedsentry.Config{
		DSN:              sentryDSN,
		Environment:      env,
		Release:          getEnv("GIT_COMMIT", "dev"),
		ServerName:       "registration-service",
		TracesSampleRate: getFloatEnv("SENTRY_TRACES_SAMPLE_RATE", 0.1),
		Debug:            env == "development",
	}); err != nil {
		log.Printf("Warning: Failed to initialize Sentry: %v", err)
	}
	defer sharedsentry.Close()

	// Integration service URL for Odoo sync
	integrationServiceURL := getEnv("INTEGRATION_SERVICE_URL", "http://localhost:8084")

	// WebSocket service URL for real-time updates
	websocketServiceURL := getEnv("WEBSOCKET_SERVICE_URL", "http://localhost:8099")

	// Initialize repositories
	registrationRepo := repository.NewRegistrationRepository(db)
	clientRepo := repository.NewClientRepository(db)

	// Initialize WebSocket client for real-time notifications
	wsClient := websocket.NewWebSocketClient(websocketServiceURL)
	log.Printf("✅ WebSocket client initialized: %s", websocketServiceURL)

	// Initialize services
	registrationService, err := services.NewRegistrationService(registrationRepo, rabbitConn, wsClient)
	if err != nil {
		log.Fatalf("Failed to create registration service: %v", err)
	}
	defer registrationService.Close()

	clientService := services.NewClientService(clientRepo)

	// Initialize handlers
	registrationHandler := handlers.NewRegistrationHandler(registrationService)
	clientHandler := handlers.NewClientHandler(clientService)

	// Initialize and start Odoo sync event consumer
	odooSyncConsumer, err := events.NewOdooSyncConsumer(rabbitConn, integrationServiceURL)
	if err != nil {
		log.Fatalf("Failed to create Odoo sync consumer: %v", err)
	}
	defer odooSyncConsumer.Close()

	// Start listening for registration.approved events
	if err := odooSyncConsumer.Start(); err != nil {
		log.Fatalf("Failed to start Odoo sync consumer: %v", err)
	}
	log.Println("✅ Odoo auto-sync enabled - listening for registration.approved events")

	// Setup router
	r := setupRouter(db, registrationHandler, clientHandler, jwtSecret)

	// PRODUCTION: Configure HTTP server with timeouts for security and reliability
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		// SECURITY: Server timeouts prevent slowloris attacks and resource exhaustion
		ReadTimeout:       15 * time.Second, // Time to read request headers and body
		ReadHeaderTimeout: 10 * time.Second, // Time to read request headers only
		WriteTimeout:      30 * time.Second, // Time to write response
		IdleTimeout:       120 * time.Second, // Keep-alive idle timeout
		MaxHeaderBytes:    1 << 20,          // 1MB max header size
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Registration Service starting on %s (with graceful shutdown support)", addr)
		log.Printf("Server timeouts: Read=%v, Write=%v, Idle=%v",
			srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// PRODUCTION: Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal %v, initiating graceful shutdown...", sig)

	// Give existing requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Registration Service stopped gracefully")
}

func setupRouter(db *sql.DB, registrationHandler *handlers.RegistrationHandler, clientHandler *handlers.ClientHandler, jwtSecret string) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// PRODUCTION: Sentry error tracking middleware (early in chain)
	r.Use(sharedsentry.Middleware())
	r.Use(sharedsentry.ErrorRecoveryMiddleware())
	r.Use(sharedsentry.PerformanceMiddleware())

	// SECURITY: Input validation middleware (before any routes)
	inputValidator := sharedmiddleware.NewInputValidator(&sharedmiddleware.ValidationConfig{
		MaxBodySize: 2 * 1024 * 1024, // 2MB limit for registrations
		AllowedMimeTypes: []string{
			"application/json",
			"application/x-www-form-urlencoded",
		},
		BlockPatterns: []string{
			`[<>]`, // No HTML brackets in company names/addresses
		},
	})
	r.Use(inputValidator.Middleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "registration-service",
			"version": "1.0.0",
		})
	})

	// PRODUCTION: Prometheus metrics endpoint for monitoring and alerting
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := r.Group("/api/v1")
	{
		// Enhanced Auth middleware - validates JWT and handles both system and tenant users
		// System users (system_admin, global_admin) have no tenant_id
		// Tenant users have tenant_id extracted from JWT claims
		api.Use(sharedmiddleware.EnhancedAuthMiddleware(jwtSecret))

		// Registration routes
		registrations := api.Group("/registrations")
		{
			registrationHandler.SetupRoutes(registrations)
		}

		// Client routes
		clients := api.Group("/clients")
		{
			clients.POST("", clientHandler.CreateClient)
			clients.GET("", clientHandler.ListClients)
			clients.GET("/:id", clientHandler.GetClient)
			clients.PUT("/:id", clientHandler.UpdateClient)
			clients.DELETE("/:id", clientHandler.DeleteClient)
		}
	}

	return r
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getFloatEnv gets a float environment variable with a default value
func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
		log.Printf("Warning: Invalid float value for %s: %s, using default: %f", key, value, defaultValue)
	}
	return defaultValue
}
