package router

import (
	"github.com/gin-gonic/gin"
)

const (
	commissionServiceURLEnvKey = "COMMISSION_SERVICE_URL"
	defaultCommissionServiceURL = "http://localhost:8085"
)

// SetupCommissionRoutes configures commission management routes
func SetupCommissionRoutes(router *gin.RouterGroup) {
	commissionServiceURL := getEnv(commissionServiceURLEnvKey, defaultCommissionServiceURL)

	// Commission CRUD operations
	router.POST("", proxyToService(commissionServiceURL, "/api/v1/commissions"))
	router.GET("", proxyToService(commissionServiceURL, "/api/v1/commissions"))
	router.GET("/:id", proxyToService(commissionServiceURL, "/api/v1/commissions/:id"))
	router.PUT("/:id", proxyToService(commissionServiceURL, "/api/v1/commissions/:id"))
	router.DELETE("/:id", proxyToService(commissionServiceURL, "/api/v1/commissions/:id"))

	// Commission calculations
	router.POST("/calculate", proxyToService(commissionServiceURL, "/api/v1/commissions/calculate"))
	router.POST("/:id/recalculate", proxyToService(commissionServiceURL, "/api/v1/commissions/:id/recalculate"))

	// Commission approval workflow
	router.POST("/:id/submit", proxyToService(commissionServiceURL, "/api/v1/commissions/:id/submit"))
	router.POST("/:id/approve", proxyToService(commissionServiceURL, "/api/v1/commissions/:id/approve"))
	router.POST("/:id/reject", proxyToService(commissionServiceURL, "/api/v1/commissions/:id/reject"))
	router.POST("/:id/pay", proxyToService(commissionServiceURL, "/api/v1/commissions/:id/pay"))

	// Commission by registration
	router.GET("/by-registration/:registration_id", proxyToService(commissionServiceURL, "/api/v1/commissions/by-registration/:registration_id"))
	router.POST("/by-registration/:registration_id/create", proxyToService(commissionServiceURL, "/api/v1/commissions/by-registration/:registration_id/create"))

	// Commission by agent/user
	router.GET("/by-agent/:agent_id", proxyToService(commissionServiceURL, "/api/v1/commissions/by-agent/:agent_id"))
	router.GET("/my-commissions", proxyToService(commissionServiceURL, "/api/v1/commissions/my-commissions"))

	// Commission rules and tiers
	router.GET("/rules", proxyToService(commissionServiceURL, "/api/v1/commissions/rules"))
	router.POST("/rules", proxyToService(commissionServiceURL, "/api/v1/commissions/rules"))
	router.PUT("/rules/:rule_id", proxyToService(commissionServiceURL, "/api/v1/commissions/rules/:rule_id"))
	router.DELETE("/rules/:rule_id", proxyToService(commissionServiceURL, "/api/v1/commissions/rules/:rule_id"))

	// Commission reports and analytics
	router.GET("/reports/summary", proxyToService(commissionServiceURL, "/api/v1/commissions/reports/summary"))
	router.GET("/reports/by-period", proxyToService(commissionServiceURL, "/api/v1/commissions/reports/by-period"))
	router.GET("/reports/by-agent", proxyToService(commissionServiceURL, "/api/v1/commissions/reports/by-agent"))
	router.GET("/reports/export", proxyToService(commissionServiceURL, "/api/v1/commissions/reports/export"))

	// Payment processing
	router.POST("/payments/batch", proxyToService(commissionServiceURL, "/api/v1/commissions/payments/batch"))
	router.GET("/payments/history", proxyToService(commissionServiceURL, "/api/v1/commissions/payments/history"))
	router.GET("/payments/:payment_id", proxyToService(commissionServiceURL, "/api/v1/commissions/payments/:payment_id"))

	// Statistics
	router.GET("/statistics", proxyToService(commissionServiceURL, "/api/v1/commissions/statistics"))
	router.GET("/statistics/pending", proxyToService(commissionServiceURL, "/api/v1/commissions/statistics/pending"))
	router.GET("/statistics/paid", proxyToService(commissionServiceURL, "/api/v1/commissions/statistics/paid"))
}
