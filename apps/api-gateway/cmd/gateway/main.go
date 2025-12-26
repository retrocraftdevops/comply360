package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/comply360/api-gateway/internal/middleware"
	"github.com/comply360/api-gateway/internal/router"
	sharedmiddleware "github.com/comply360/shared/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	dbURL := getEnv("DATABASE_URL", "postgresql://comply360_app_user:comply360_app_secure_pass@localhost:5432/comply360_app?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379/0")
	port := getEnv("API_GATEWAY_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "dev_secret_key_change_in_production")

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Connect to Redis
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// Setup router
	r := setupRouter(db, redisClient, jwtSecret)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("API Gateway starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(db *sql.DB, redisClient *redis.Client, jwtSecret string) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:5176"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Tenant-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,  // Changed to false for direct API calls
		MaxAge:           12 * time.Hour,
	}))

	// Global middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())

	// Health check (no tenant required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "api-gateway",
			"version": "1.0.0",
		})
	})

	// API routes with tenant middleware
	api := r.Group("/api")
	{
		// Tenant middleware - extracts tenant context
		api.Use(sharedmiddleware.TenantMiddleware(db))

		// Rate limiting middleware
		api.Use(middleware.RateLimiter(redisClient))

		// Route to backend services
		v1 := api.Group("/v1")
		{
			// Auth service routes (public, no auth required)
			auth := v1.Group("/auth")
			{
				router.SetupAuthRoutes(auth)
			}

			// Tenant service routes (admin only)
			tenants := v1.Group("/tenants")
			{
				router.SetupTenantRoutes(tenants)
			}

			// Registration service routes (authenticated)
			registrations := v1.Group("/registrations")
			registrations.Use(sharedmiddleware.AuthMiddleware(jwtSecret))
			{
				router.SetupRegistrationRoutes(registrations)
			}

			// Document service routes (authenticated)
			documents := v1.Group("/documents")
			documents.Use(sharedmiddleware.AuthMiddleware(jwtSecret))
			{
				router.SetupDocumentRoutes(documents)
			}

			// Commission service routes (authenticated)
			commissions := v1.Group("/commissions")
			commissions.Use(sharedmiddleware.AuthMiddleware(jwtSecret))
			{
				router.SetupCommissionRoutes(commissions)
			}

			// Integration service routes (internal only)
			integration := v1.Group("/integration")
			{
				router.SetupIntegrationRoutes(integration)
			}
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
