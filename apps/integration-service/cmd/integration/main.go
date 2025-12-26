package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/comply360/integration-service/internal/adapters"
	"github.com/comply360/integration-service/internal/consumers"
	"github.com/comply360/integration-service/internal/handlers"
	"github.com/comply360/integration-service/internal/services"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration from environment
	port := getEnv("INTEGRATION_SERVICE_PORT", "8086")
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://comply360:dev_password@localhost:5672/")

	// Odoo configuration
	odooConfig := &adapters.OdooConfig{
		URL:      getEnv("ODOO_URL", "http://localhost:8069"),
		Database: getEnv("ODOO_DATABASE", "comply360_odoo"),
		Username: getEnv("ODOO_USERNAME", "admin"),
		Password: getEnv("ODOO_PASSWORD", "admin"),
	}

	// Initialize Odoo client
	odooClient, err := adapters.NewOdooClient(odooConfig)
	if err != nil {
		log.Printf("Warning: Failed to connect to Odoo: %v", err)
		log.Println("Integration service will start but Odoo features will be unavailable")
		// Don't fatal - allow service to start for health checks
	} else {
		log.Println("Successfully connected to Odoo")
	}

	// Initialize services
	odooService := services.NewOdooService(odooClient)

	// Initialize handlers
	odooHandler := handlers.NewOdooHandler(odooService)

	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to RabbitMQ: %v", err)
		log.Println("Integration service will start but event consumption will be unavailable")
	} else {
		log.Println("Connected to RabbitMQ")

		// Initialize and start Odoo sync consumer
		odooSyncConsumer, err := consumers.NewOdooSyncConsumer(rabbitConn, odooService)
		if err != nil {
			log.Printf("Warning: Failed to create Odoo sync consumer: %v", err)
		} else {
			if err := odooSyncConsumer.Start(); err != nil {
				log.Printf("Warning: Failed to start Odoo sync consumer: %v", err)
			} else {
				defer odooSyncConsumer.Close()
			}
		}
		defer rabbitConn.Close()
	}

	// Setup router
	r := setupRouter(odooHandler)

	// Start HTTP server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Integration Service starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Integration Service...")
}

func setupRouter(odooHandler *handlers.OdooHandler) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "integration-service",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api/v1/integration")
	{
		// Odoo integration routes
		odoo := api.Group("/odoo")
		{
			// Lead management
			odoo.POST("/leads", odooHandler.CreateLead)
			odoo.GET("/leads/:id", odooHandler.GetLead)
			odoo.PUT("/leads/:id", odooHandler.UpdateLead)
			odoo.POST("/leads/:id/convert", odooHandler.ConvertLeadToCustomer)

			// Customer management
			// odoo.GET("/customers/:id", odooHandler.GetCustomer)
			// odoo.PUT("/customers/:id", odooHandler.UpdateCustomer)

			// Invoice management
			odoo.POST("/invoices", odooHandler.CreateInvoice)
			// odoo.GET("/invoices/:id", odooHandler.GetInvoice)

			// Commission management
			odoo.POST("/commissions", odooHandler.CreateCommission)
			// odoo.GET("/commissions/:id", odooHandler.GetCommission)

			// Sync operations
			odoo.POST("/sync/registration/:registration_id", odooHandler.SyncRegistration)
			odoo.POST("/sync/commission/:commission_id", odooHandler.SyncCommission)

			// Connection management
			odoo.GET("/status", odooHandler.GetStatus)
			odoo.POST("/test-connection", odooHandler.TestConnection)
		}

		// CIPC integration routes (placeholder)
		cipc := api.Group("/cipc")
		{
			cipc.POST("/search", placeholderHandler("CIPC search"))
			cipc.POST("/verify", placeholderHandler("CIPC verify"))
			cipc.GET("/company/:registration_number", placeholderHandler("CIPC company lookup"))
			cipc.GET("/status", placeholderHandler("CIPC status"))
		}

		// SARS integration routes (placeholder)
		sars := api.Group("/sars")
		{
			sars.POST("/verify-vat", placeholderHandler("SARS VAT verification"))
			sars.POST("/verify-tax-number", placeholderHandler("SARS tax number verification"))
			sars.GET("/status", placeholderHandler("SARS status"))
		}

		// Payment gateway routes (placeholder)
		payments := api.Group("/payments")
		{
			// Stripe
			payments.POST("/stripe/create-payment-intent", placeholderHandler("Stripe payment intent"))
			payments.POST("/stripe/webhook", placeholderHandler("Stripe webhook"))

			// PayFast
			payments.POST("/payfast/create-payment", placeholderHandler("PayFast payment"))
			payments.POST("/payfast/webhook", placeholderHandler("PayFast webhook"))

			payments.GET("/status/:payment_id", placeholderHandler("Payment status"))
		}

		// Email service routes (placeholder)
		email := api.Group("/email")
		{
			email.POST("/send", placeholderHandler("Send email"))
			email.POST("/send-template", placeholderHandler("Send template email"))
			email.GET("/templates", placeholderHandler("List email templates"))
		}

		// SMS service routes (placeholder)
		sms := api.Group("/sms")
		{
			sms.POST("/send", placeholderHandler("Send SMS"))
			sms.GET("/status/:message_id", placeholderHandler("SMS status"))
		}

		// Webhook management (placeholder)
		webhooks := api.Group("/webhooks")
		{
			webhooks.GET("", placeholderHandler("List webhooks"))
			webhooks.POST("", placeholderHandler("Create webhook"))
			webhooks.GET("/:id", placeholderHandler("Get webhook"))
			webhooks.PUT("/:id", placeholderHandler("Update webhook"))
			webhooks.DELETE("/:id", placeholderHandler("Delete webhook"))
			webhooks.POST("/:id/test", placeholderHandler("Test webhook"))
		}

		// Overall integration status
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"service": "integration-service",
				"integrations": gin.H{
					"odoo":    "connected",
					"cipc":    "not_configured",
					"sars":    "not_configured",
					"stripe":  "not_configured",
					"payfast": "not_configured",
				},
			})
		})
	}

	return r
}

func placeholderHandler(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(501, gin.H{
			"message": fmt.Sprintf("%s not implemented yet", feature),
		})
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
