package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/comply360/document-service/internal/handlers"
	"github.com/comply360/document-service/internal/repository"
	"github.com/comply360/document-service/internal/services"
	sharedmiddleware "github.com/comply360/shared/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	dbURL := getEnv("DATABASE_URL", "postgresql://comply360_app_user:comply360_app_secure_pass@localhost:5432/comply360_app?sslmode=disable")
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://comply360:dev_password@localhost:5672/")
	port := getEnv("DOCUMENT_SERVICE_PORT", "8084")
	jwtSecret := getEnv("JWT_SECRET", "dev_secret_key_change_in_production")

	// MinIO configuration
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getEnv("MINIO_BUCKET", "comply360-documents")
	minioUseSSL := getEnv("MINIO_USE_SSL", "false") == "true"

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

	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	// Initialize MinIO client
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: minioUseSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}
	log.Println("Connected to MinIO")

	// Initialize repository
	repo := repository.NewDocumentRepository(db)

	// Initialize service
	service, err := services.NewDocumentService(repo, minioClient, minioBucket, rabbitConn)
	if err != nil {
		log.Fatalf("Failed to create document service: %v", err)
	}
	defer service.Close()

	// Initialize handler
	handler := handlers.NewDocumentHandler(service)

	// Setup router
	r := setupRouter(db, handler, jwtSecret)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Document Service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(db *sql.DB, handler *handlers.DocumentHandler, jwtSecret string) *gin.Engine {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "document-service",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Tenant middleware - extracts tenant context
		api.Use(sharedmiddleware.TenantMiddleware(db))

		// Auth middleware - validates JWT
		api.Use(sharedmiddleware.AuthMiddleware(jwtSecret))

		// Document routes
		documents := api.Group("/documents")
		{
			handler.SetupRoutes(documents)
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
