package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig returns default CORS configuration for development
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",  // Vite default (SvelteKit)
			"http://localhost:4173",  // Vite preview
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:4173",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// CORS returns a CORS middleware configured for the application
func CORS(config CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Tenant-ID",
			"X-Request-ID",
			"X-CSRF-Token",
			"Origin",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	}

	return cors.New(corsConfig)
}

// CORSForProduction returns CORS middleware for production with specific origins
func CORSForProduction(allowedOrigins []string) gin.HandlerFunc {
	return CORS(CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
