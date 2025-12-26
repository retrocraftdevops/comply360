package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// proxyToService creates a generic handler that proxies requests to a backend service
func proxyToService(baseURL, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build target URL
		targetURL := baseURL + path

		// Replace path parameters
		for _, param := range c.Params {
			placeholder := ":" + param.Key
			targetURL = strings.Replace(targetURL, placeholder, param.Value, 1)
		}

		// Add query parameters
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create new request
		req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": "PROXY_ERROR",
				"message": "Failed to create proxy request",
			})
			return
		}

		// Copy headers from original request
		for key, values := range c.Request.Header {
			// Skip certain headers
			if shouldSkipHeader(key) {
				continue
			}
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Add tenant context header if available
		if tenantID, exists := c.Get("tenant_id"); exists {
			req.Header.Set("X-Tenant-ID", fmt.Sprintf("%v", tenantID))
		}

		// Add request ID if available
		if requestID, exists := c.Get("request_id"); exists {
			req.Header.Set("X-Request-ID", fmt.Sprintf("%v", requestID))
		}

		// Forward request to backend service
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"code": "SERVICE_UNAVAILABLE",
				"message": fmt.Sprintf("Backend service unavailable: %s", baseURL),
			})
			return
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": "PROXY_ERROR",
				"message": "Failed to read backend service response",
			})
			return
		}

		// Copy response headers
		for key, values := range resp.Header {
			if shouldSkipHeader(key) {
				continue
			}
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Parse JSON response if possible
		var jsonResponse interface{}
		if err := json.Unmarshal(respBody, &jsonResponse); err == nil {
			c.JSON(resp.StatusCode, jsonResponse)
		} else {
			// Non-JSON response (could be file download, etc.)
			c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
		}
	}
}

// shouldSkipHeader determines if a header should be skipped when proxying
func shouldSkipHeader(header string) bool {
	header = strings.ToLower(header)
	skipHeaders := []string{
		"content-length",
		"transfer-encoding",
		"connection",
		"keep-alive",
		"proxy-authenticate",
		"proxy-authorization",
		"te",
		"trailers",
		"upgrade",
	}

	for _, skip := range skipHeaders {
		if header == skip {
			return true
		}
	}
	return false
}

// getEnv retrieves an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
