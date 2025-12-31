package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/comply360/notification-service/internal/consumers"
	"github.com/comply360/notification-service/internal/handlers"
	"github.com/comply360/notification-service/internal/repository"
	"github.com/comply360/notification-service/internal/services"
	sharedmiddleware "github.com/comply360/shared/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	// SECURITY: Default to SSL required for production-grade security
	dbURL := getEnv("DATABASE_URL", "postgres://comply360:dev_password@localhost:5432/comply360_dev?sslmode=require")
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://comply360:dev_password@localhost:5672/")
	port := getEnv("NOTIFICATION_SERVICE_PORT", "8087")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Email configuration
	smtpHost := getEnv("SMTP_HOST", "smtp.example.com")
	smtpPort := getEnvInt("SMTP_PORT", 587)
	smtpUsername := getEnv("SMTP_USERNAME", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	fromEmail := getEnv("FROM_EMAIL", "noreply@comply360.africa")
	fromName := getEnv("FROM_NAME", "Comply360")

	// SMS configuration
	smsAPIKey := getEnv("SMS_API_KEY", "")
	smsAPISecret := getEnv("SMS_API_SECRET", "")
	smsFromNumber := getEnv("SMS_FROM_NUMBER", "")
	smsProvider := getEnv("SMS_PROVIDER", "twilio")

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

	// Create sqlx wrapper
	dbx := sqlx.NewDb(db, "postgres")

	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	// Initialize repositories
	notificationRepo := repository.NewNotificationRepository(dbx)

	// Initialize email service
	emailService := services.NewEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, fromName)

	// Initialize SMS service
	smsService := services.NewSMSService(smsAPIKey, smsAPISecret, smsFromNumber, smsProvider)

	// Initialize notification service
	notificationService := services.NewNotificationService(notificationRepo)

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// Initialize event consumer
	eventConsumer, err := consumers.NewEventConsumer(rabbitConn, emailService, smsService)
	if err != nil {
		log.Fatalf("Failed to create event consumer: %v", err)
	}
	defer eventConsumer.Close()

	// Start consuming events
	if err := eventConsumer.Start(); err != nil {
		log.Fatalf("Failed to start event consumer: %v", err)
	}

	// Setup HTTP server
	r := setupRouter(emailHandler, notificationHandler, jwtSecret)

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Notification Service starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Notification Service...")
}

func setupRouter(emailHandler *handlers.EmailHandler, notificationHandler *handlers.NotificationHandler, jwtSecret string) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "notification-service",
			"version": "1.0.0",
		})
	})

	// PRODUCTION: Prometheus metrics endpoint for monitoring and alerting
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := r.Group("/api/v1/notifications")
	{
		// Enhanced Auth middleware - validates JWT and handles both system and tenant users
		// System users (system_admin, global_admin) may have empty tenant_id
		if jwtSecret != "" {
			api.Use(sharedmiddleware.EnhancedAuthMiddleware(jwtSecret))
		}
		// Email routes
		email := api.Group("/email")
		{
			emailHandler.SetupRoutes(email)
		}

		// SMS routes (placeholder for future implementation)
		sms := api.Group("/sms")
		{
			sms.POST("/send", func(c *gin.Context) {
				c.JSON(501, gin.H{
					"message": "SMS functionality not yet implemented",
				})
			})
		}

		// In-app notification routes
		notificationHandler.SetupRoutes(api)
	}

	return r
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
