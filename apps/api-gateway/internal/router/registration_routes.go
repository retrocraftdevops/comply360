package router

import (
	"github.com/gin-gonic/gin"
)

const (
	registrationServiceURLEnvKey = "REGISTRATION_SERVICE_URL"
	defaultRegistrationServiceURL = "http://localhost:8081"
)

// SetupRegistrationRoutes configures registration management routes
func SetupRegistrationRoutes(router *gin.RouterGroup) {
	registrationServiceURL := getEnv(registrationServiceURLEnvKey, defaultRegistrationServiceURL)

	// Registration CRUD operations
	router.POST("", proxyToService(registrationServiceURL, "/api/v1/registrations"))
	router.GET("", proxyToService(registrationServiceURL, "/api/v1/registrations"))
	router.GET("/:id", proxyToService(registrationServiceURL, "/api/v1/registrations/:id"))
	router.PUT("/:id", proxyToService(registrationServiceURL, "/api/v1/registrations/:id"))
	router.DELETE("/:id", proxyToService(registrationServiceURL, "/api/v1/registrations/:id"))

	// Registration workflow actions
	router.POST("/:id/submit", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/submit"))
	router.POST("/:id/approve", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/approve"))
	router.POST("/:id/reject", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/reject"))
	router.POST("/:id/cancel", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/cancel"))

	// Document verification
	router.POST("/:id/verify-documents", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/verify-documents"))

	// CIPC integration
	router.POST("/:id/cipc-search", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/cipc-search"))
	router.POST("/:id/cipc-verify", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/cipc-verify"))

	// Registration history
	router.GET("/:id/history", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/history"))
	router.GET("/:id/audit", proxyToService(registrationServiceURL, "/api/v1/registrations/:id/audit"))

	// Statistics and reports
	router.GET("/statistics", proxyToService(registrationServiceURL, "/api/v1/registrations/statistics"))
	router.GET("/export", proxyToService(registrationServiceURL, "/api/v1/registrations/export"))
}
