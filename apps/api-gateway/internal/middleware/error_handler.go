package middleware

import (
	"log"
	"net/http"

	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware catches panics and handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover from panics
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)

				log.Printf("[PANIC] Request ID: %s | Error: %v", requestID, err)

				// Return internal server error
				c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
					errors.ErrInternalServer,
					"An unexpected error occurred",
					map[string]interface{}{
						"request_id": requestID,
					},
				))

				c.Abort()
			}
		}()

		// Process request
		c.Next()

		// Handle errors added during request processing
		if len(c.Errors) > 0 {
			requestID := GetRequestID(c)

			// Get the last error (most recent)
			err := c.Errors.Last()

			// Check if it's an APIError
			if apiErr, ok := err.Err.(*errors.APIError); ok {
				// Determine HTTP status code from error code
				statusCode := getHTTPStatusFromErrorCode(apiErr.Code)
				c.JSON(statusCode, apiErr)
				return
			}

			// Handle generic errors
			log.Printf("[ERROR] Request ID: %s | Error: %v", requestID, err.Err)

			c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
				errors.ErrInternalServer,
				err.Err.Error(),
				map[string]interface{}{
					"request_id": requestID,
				},
			))
		}
	}
}

// getHTTPStatusFromErrorCode maps error codes to HTTP status codes
func getHTTPStatusFromErrorCode(code string) int {
	switch code {
	case errors.ErrInvalidInput:
		return http.StatusBadRequest
	case errors.ErrUnauthorized, errors.ErrInvalidToken, errors.ErrTokenExpired:
		return http.StatusUnauthorized
	case errors.ErrForbidden, errors.ErrInsufficientPermissions, errors.ErrTenantSuspended:
		return http.StatusForbidden
	case errors.ErrNotFound, errors.ErrTenantNotFound:
		return http.StatusNotFound
	case errors.ErrConflict:
		return http.StatusConflict
	case errors.ErrRateLimitExceeded:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
