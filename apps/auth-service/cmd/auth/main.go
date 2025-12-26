package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/comply360/auth-service/internal/handlers"
	"github.com/comply360/auth-service/internal/repository"
	"github.com/comply360/auth-service/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration from environment
	dbURL := getEnv("DATABASE_URL", "postgres://comply360:dev_password@localhost:5432/comply360_dev?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379/0")
	port := getEnv("AUTH_SERVICE_PORT", "8081")
	jwtSecret := getEnv("JWT_SECRET", "changeme-change-this-in-production")

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

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, redisClient, jwtSecret)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	r := setupRouter(authHandler)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Auth Service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "auth-service",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api/v1/auth")
	{
		// Public authentication endpoints
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)
		api.POST("/forgot-password", authHandler.ForgotPassword)
		api.POST("/reset-password", authHandler.ResetPassword)
		api.POST("/verify-email", authHandler.VerifyEmail)
		api.POST("/resend-verification", authHandler.ResendVerification)

		// OAuth endpoints
		oauth := api.Group("/oauth")
		{
			oauth.GET("/:provider", authHandler.OAuthLogin)
			oauth.GET("/:provider/callback", authHandler.OAuthCallback)
		}

		// MFA endpoints (require authentication)
		mfa := api.Group("/mfa")
		{
			mfa.POST("/setup", authHandler.SetupMFA)
			mfa.POST("/verify", authHandler.VerifyMFA)
			mfa.POST("/disable", authHandler.DisableMFA)
		}

		// Password management (authenticated)
		api.POST("/change-password", authHandler.ChangePassword)

		// Logout
		api.POST("/logout", authHandler.Logout)

		// User profile (authenticated)
		api.GET("/me", authHandler.GetProfile)
		api.PUT("/me", authHandler.UpdateProfile)
	}

	return r
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
