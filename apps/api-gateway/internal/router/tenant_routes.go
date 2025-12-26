package router

import (
	"github.com/gin-gonic/gin"
)

const (
	tenantServiceURLEnvKey = "TENANT_SERVICE_URL"
	defaultTenantServiceURL = "http://localhost:8082"
)

// SetupTenantRoutes configures tenant management routes
// These routes are typically admin-only
func SetupTenantRoutes(router *gin.RouterGroup) {
	tenantServiceURL := getEnv(tenantServiceURLEnvKey, defaultTenantServiceURL)

	// Tenant CRUD operations
	router.POST("", proxyToService(tenantServiceURL, "/api/v1/tenants"))
	router.GET("", proxyToService(tenantServiceURL, "/api/v1/tenants"))
	router.GET("/:id", proxyToService(tenantServiceURL, "/api/v1/tenants/:id"))
	router.PUT("/:id", proxyToService(tenantServiceURL, "/api/v1/tenants/:id"))
	router.DELETE("/:id", proxyToService(tenantServiceURL, "/api/v1/tenants/:id"))

	// Tenant provisioning
	router.POST("/:id/provision", proxyToService(tenantServiceURL, "/api/v1/tenants/:id/provision"))

	// Tenant settings
	router.GET("/:id/settings", proxyToService(tenantServiceURL, "/api/v1/tenants/:id/settings"))
	router.PUT("/:id/settings", proxyToService(tenantServiceURL, "/api/v1/tenants/:id/settings"))

	// Tenant users
	router.GET("/:id/users", proxyToService(tenantServiceURL, "/api/v1/tenants/:id/users"))
	router.GET("/:id/statistics", proxyToService(tenantServiceURL, "/api/v1/tenants/:id/statistics"))
}
