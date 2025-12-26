package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/comply360/tenant-service/internal/handlers"
	"github.com/comply360/tenant-service/internal/repository"
	"github.com/comply360/tenant-service/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	dbURL := getEnv("DATABASE_URL", "postgres://comply360_user:comply360_password_dev@localhost:5432/comply360_db?sslmode=disable")
	port := getEnv("TENANT_SERVICE_PORT", "8082")

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Initialize repository
	tenantRepo := repository.NewTenantRepository(db)

	// Initialize service
	tenantService := services.NewTenantService(tenantRepo, db)

	// Initialize handlers
	tenantHandler := handlers.NewTenantHandler(tenantService)

	// Setup Gin router
	router := setupRouter(tenantHandler)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Tenant Service starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(tenantHandler *handlers.TenantHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "tenant-service",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Tenant management routes (admin only)
		tenants := v1.Group("/tenants")
		{
			tenants.POST("", tenantHandler.CreateTenant)
			tenants.GET("", tenantHandler.ListTenants)
			tenants.GET("/:id", tenantHandler.GetTenant)
			tenants.PUT("/:id", tenantHandler.UpdateTenant)
			tenants.DELETE("/:id", tenantHandler.DeleteTenant)
			tenants.POST("/:id/provision", tenantHandler.ProvisionTenant)
		}
	}

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
