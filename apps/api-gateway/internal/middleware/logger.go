package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware logs HTTP requests with structured information
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get request ID
		requestID := GetRequestID(c)

		// Get tenant ID if available
		tenantID := ""
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = tid.(string)
		}

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Get error if any
		errorMessage := ""
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Log request
		log.Printf(
			"[%s] %s | %3d | %13v | %15s | %-7s %s | Tenant: %s | Error: %s",
			requestID,
			startTime.Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			tenantID,
			errorMessage,
		)

		// Log slow requests (>1s)
		if latency > time.Second {
			log.Printf("[SLOW REQUEST] %s | %s %s | %v | Tenant: %s",
				requestID,
				c.Request.Method,
				c.Request.URL.Path,
				latency,
				tenantID,
			)
		}
	}
}
