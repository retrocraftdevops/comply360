package router

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	sharedmiddleware "github.com/comply360/shared/middleware"
)

// SetupAdminRoutes configures all admin-related routes
func SetupAdminRoutes(router *gin.Engine, db *sql.DB) {
	// Get auth service URL
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8081")

	// Admin routes group - requires authentication and admin/manager role
	admin := router.Group("/api/v1/admin")
	admin.Use(sharedmiddleware.EnhancedAuthMiddleware(os.Getenv("JWT_SECRET")))
	// Don't require tenant middleware for system admins
	// Require admin or manager role for all admin routes
	// Enhanced auth includes role hierarchy support
	// system_admin has god-mode and bypasses all restrictions
	admin.Use(sharedmiddleware.RequireAnyRole("system_admin", "global_admin", "tenant_admin", "tenant_manager"))
	// TODO: Implement AuditMiddleware for logging admin actions
	// admin.Use(AuditMiddleware(db))

	// User Management Routes
	userRoutes := admin.Group("/users")
	{
		// View operations - requires users.view permission (handled by auth-service)
		userRoutes.GET("", proxyToService(authServiceURL, "/api/v1/users"))
		userRoutes.GET("/:id", proxyToService(authServiceURL, "/api/v1/users/:id"))
		userRoutes.GET("/:id/effective-permissions", proxyToService(authServiceURL, "/api/v1/users/:id/effective-permissions"))

		// Create/Update/Delete operations - permissions checked by auth-service
		// users.create, users.edit, users.delete permissions required
		userRoutes.POST("", proxyToService(authServiceURL, "/api/v1/users"))
		userRoutes.PUT("/:id", proxyToService(authServiceURL, "/api/v1/users/:id"))
		userRoutes.DELETE("/:id", proxyToService(authServiceURL, "/api/v1/users/:id"))

		// User management actions - permissions checked by auth-service
		userRoutes.POST("/:id/activate", proxyToService(authServiceURL, "/api/v1/users/:id/activate"))
		userRoutes.POST("/:id/deactivate", proxyToService(authServiceURL, "/api/v1/users/:id/deactivate"))
		userRoutes.POST("/:id/unlock", proxyToService(authServiceURL, "/api/v1/users/:id/unlock"))
	}

	// Role Management Routes
	roleRoutes := admin.Group("/roles")
	{
		// View roles and permissions
		roleRoutes.GET("", proxyToService(authServiceURL, "/api/v1/roles"))
		roleRoutes.GET("/:role/permissions", proxyToService(authServiceURL, "/api/v1/roles/:role/permissions"))

		// User role management - requires users.manage_roles permission (checked by auth-service)
		roleRoutes.GET("/users/:id/roles", proxyToService(authServiceURL, "/api/v1/users/:id/roles"))
		roleRoutes.POST("/users/:id/roles", proxyToService(authServiceURL, "/api/v1/users/:id/roles"))
		roleRoutes.DELETE("/users/:id/roles/:role", proxyToService(authServiceURL, "/api/v1/users/:id/roles/:role"))
		roleRoutes.GET("/users/:id/effective-permissions", proxyToService(authServiceURL, "/api/v1/users/:id/effective-permissions"))
	}

	// Feature Management Routes
	featureRoutes := admin.Group("/features")
	{
		// Anyone with admin access can view features
		featureRoutes.GET("", proxyToService(authServiceURL, "/api/v1/features"))
		featureRoutes.GET("/enabled", proxyToService(authServiceURL, "/api/v1/features/enabled"))
		featureRoutes.GET("/:code/check", proxyToService(authServiceURL, "/api/v1/features/:code/check"))

		// Feature modification routes - only tenant_admin allowed
		// Role level 2 or better (tenant_admin, global_admin, system_admin)
		featureModify := featureRoutes.Group("")
		featureModify.Use(sharedmiddleware.RequireRoleLevel(2))
		{
			featureModify.POST("/:code/enable", proxyToService(authServiceURL, "/api/v1/features/:code/enable"))
			featureModify.POST("/:code/disable", proxyToService(authServiceURL, "/api/v1/features/:code/disable"))
		}
	}

	// Plan features (public - no auth required)
	router.GET("/api/v1/plans/features", proxyToService(authServiceURL, "/api/v1/plans/features"))

	// Audit Log Routes - requires admin.audit_logs permission (checked by auth-service)
	auditRoutes := admin.Group("/audit-logs")
	{
		auditRoutes.GET("", proxyToService(authServiceURL, "/api/v1/audit-logs"))
		auditRoutes.GET("/:id", proxyToService(authServiceURL, "/api/v1/audit-logs/:id"))
		auditRoutes.GET("/user-activity", proxyToService(authServiceURL, "/api/v1/audit-logs/user-activity"))
		auditRoutes.GET("/stats", proxyToService(authServiceURL, "/api/v1/audit-logs/stats"))
		auditRoutes.GET("/export", proxyToService(authServiceURL, "/api/v1/audit-logs/export"))
	}

	// System Health Routes - requires admin.system_health permission (checked by auth-service)
	systemRoutes := admin.Group("/system")
	{
		systemRoutes.GET("/health", proxyToService(authServiceURL, "/api/v1/system/health"))
		systemRoutes.GET("/services", proxyToService(authServiceURL, "/api/v1/system/services"))
		systemRoutes.GET("/database", proxyToService(authServiceURL, "/api/v1/system/database"))
		systemRoutes.GET("/metrics", proxyToService(authServiceURL, "/api/v1/system/metrics"))
		systemRoutes.GET("/usage", proxyToService(authServiceURL, "/api/v1/system/usage"))
	}

	// Permissions info routes (for frontend to build permission-based UI)
	permissionRoutes := admin.Group("/permissions")
	{
		permissionRoutes.GET("/me", func(c *gin.Context) {
			// This endpoint returns current user's permissions
			// Will be implemented to call permission service
			c.JSON(http.StatusOK, gin.H{
				"message": "Get current user permissions - to be implemented",
			})
		})
	}
}
