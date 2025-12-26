package router

import (
	"github.com/gin-gonic/gin"
)

const (
	documentServiceURLEnvKey = "DOCUMENT_SERVICE_URL"
	defaultDocumentServiceURL = "http://localhost:8084"
)

// SetupDocumentRoutes configures document management routes
func SetupDocumentRoutes(router *gin.RouterGroup) {
	documentServiceURL := getEnv(documentServiceURLEnvKey, defaultDocumentServiceURL)

	// Document CRUD operations
	router.POST("", proxyToService(documentServiceURL, "/api/v1/documents"))
	router.GET("", proxyToService(documentServiceURL, "/api/v1/documents"))
	router.GET("/:id", proxyToService(documentServiceURL, "/api/v1/documents/:id"))
	router.PUT("/:id", proxyToService(documentServiceURL, "/api/v1/documents/:id"))
	router.DELETE("/:id", proxyToService(documentServiceURL, "/api/v1/documents/:id"))

	// Document upload and download
	router.POST("/upload", proxyToService(documentServiceURL, "/api/v1/documents/upload"))
	router.GET("/:id/download", proxyToService(documentServiceURL, "/api/v1/documents/:id/download"))
	router.GET("/:id/preview", proxyToService(documentServiceURL, "/api/v1/documents/:id/preview"))

	// Document verification
	router.POST("/:id/verify", proxyToService(documentServiceURL, "/api/v1/documents/:id/verify"))
	router.POST("/:id/ai-analyze", proxyToService(documentServiceURL, "/api/v1/documents/:id/ai-analyze"))
	router.GET("/:id/verification-status", proxyToService(documentServiceURL, "/api/v1/documents/:id/verification-status"))

	// Document versions
	router.GET("/:id/versions", proxyToService(documentServiceURL, "/api/v1/documents/:id/versions"))
	router.POST("/:id/versions", proxyToService(documentServiceURL, "/api/v1/documents/:id/versions"))
	router.GET("/:id/versions/:version_id", proxyToService(documentServiceURL, "/api/v1/documents/:id/versions/:version_id"))

	// Document sharing and permissions
	router.POST("/:id/share", proxyToService(documentServiceURL, "/api/v1/documents/:id/share"))
	router.GET("/:id/permissions", proxyToService(documentServiceURL, "/api/v1/documents/:id/permissions"))
	router.PUT("/:id/permissions", proxyToService(documentServiceURL, "/api/v1/documents/:id/permissions"))

	// Bulk operations
	router.POST("/bulk-upload", proxyToService(documentServiceURL, "/api/v1/documents/bulk-upload"))
	router.POST("/bulk-delete", proxyToService(documentServiceURL, "/api/v1/documents/bulk-delete"))
	router.POST("/bulk-verify", proxyToService(documentServiceURL, "/api/v1/documents/bulk-verify"))

	// Document templates
	router.GET("/templates", proxyToService(documentServiceURL, "/api/v1/documents/templates"))
	router.POST("/templates", proxyToService(documentServiceURL, "/api/v1/documents/templates"))
	router.GET("/templates/:template_id", proxyToService(documentServiceURL, "/api/v1/documents/templates/:template_id"))

	// Storage statistics
	router.GET("/storage/usage", proxyToService(documentServiceURL, "/api/v1/documents/storage/usage"))
	router.GET("/storage/quota", proxyToService(documentServiceURL, "/api/v1/documents/storage/quota"))
}
