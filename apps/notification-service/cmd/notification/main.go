package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/comply360/notification-service/internal/consumers"
	"github.com/comply360/notification-service/internal/services"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://comply360:dev_password@localhost:5672/")
	port := getEnv("NOTIFICATION_SERVICE_PORT", "8087")

	// Email configuration
	smtpHost := getEnv("SMTP_HOST", "smtp.example.com")
	smtpPort := getEnvInt("SMTP_PORT", 587)
	smtpUsername := getEnv("SMTP_USERNAME", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	fromEmail := getEnv("FROM_EMAIL", "noreply@comply360.com")
	fromName := getEnv("FROM_NAME", "Comply360")

	// SMS configuration
	smsAPIKey := getEnv("SMS_API_KEY", "")
	smsAPISecret := getEnv("SMS_API_SECRET", "")
	smsFromNumber := getEnv("SMS_FROM_NUMBER", "")
	smsProvider := getEnv("SMS_PROVIDER", "twilio")

	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	// Initialize email service
	emailService := services.NewEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, fromName)

	// Initialize SMS service
	smsService := services.NewSMSService(smsAPIKey, smsAPISecret, smsFromNumber, smsProvider)

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

	// Setup HTTP server for health checks
	r := setupRouter()

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

func setupRouter() *gin.Engine {
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
