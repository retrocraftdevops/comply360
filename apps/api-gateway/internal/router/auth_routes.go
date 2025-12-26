package router

import (
	"github.com/gin-gonic/gin"
)

const (
	authServiceURLEnvKey = "AUTH_SERVICE_URL"
	defaultAuthServiceURL = "http://localhost:8081"
)

// SetupAuthRoutes configures authentication routes
func SetupAuthRoutes(router *gin.RouterGroup) {
	authServiceURL := getEnv(authServiceURLEnvKey, defaultAuthServiceURL)

	// Public authentication endpoints (no auth required)
	router.POST("/register", proxyToService(authServiceURL, "/api/v1/auth/register"))
	router.POST("/login", proxyToService(authServiceURL, "/api/v1/auth/login"))
	router.POST("/refresh", proxyToService(authServiceURL, "/api/v1/auth/refresh"))
	router.POST("/forgot-password", proxyToService(authServiceURL, "/api/v1/auth/forgot-password"))
	router.POST("/reset-password", proxyToService(authServiceURL, "/api/v1/auth/reset-password"))
	router.POST("/verify-email", proxyToService(authServiceURL, "/api/v1/auth/verify-email"))
	router.POST("/resend-verification", proxyToService(authServiceURL, "/api/v1/auth/resend-verification"))

	// OAuth endpoints
	oauth := router.Group("/oauth")
	{
		oauth.GET("/:provider", proxyToService(authServiceURL, "/api/v1/auth/oauth/:provider"))
		oauth.GET("/:provider/callback", proxyToService(authServiceURL, "/api/v1/auth/oauth/:provider/callback"))
	}

	// MFA endpoints (require initial authentication)
	mfa := router.Group("/mfa")
	{
		mfa.POST("/setup", proxyToService(authServiceURL, "/api/v1/auth/mfa/setup"))
		mfa.POST("/verify", proxyToService(authServiceURL, "/api/v1/auth/mfa/verify"))
		mfa.POST("/disable", proxyToService(authServiceURL, "/api/v1/auth/mfa/disable"))
	}

	// Password management (authenticated)
	router.POST("/change-password", proxyToService(authServiceURL, "/api/v1/auth/change-password"))

	// Logout
	router.POST("/logout", proxyToService(authServiceURL, "/api/v1/auth/logout"))

	// User profile (authenticated)
	router.GET("/me", proxyToService(authServiceURL, "/api/v1/auth/me"))
	router.PUT("/me", proxyToService(authServiceURL, "/api/v1/auth/me"))
}
