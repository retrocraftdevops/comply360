package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey is the context key for request ID
	RequestIDKey = "request_id"
)

// RequestID middleware generates a unique ID for each request
// and adds it to the response headers and context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request already has an ID
		requestID := c.GetHeader(RequestIDHeader)

		// Generate new ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set in context for use by other middleware/handlers
		c.Set(RequestIDKey, requestID)

		// Set in response headers
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}
