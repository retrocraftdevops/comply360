package router

import (
	"github.com/gin-gonic/gin"
)

const (
	integrationServiceURLEnvKey = "INTEGRATION_SERVICE_URL"
	defaultIntegrationServiceURL = "http://localhost:8086"
)

// SetupIntegrationRoutes configures integration service routes
// These routes are typically internal-only for service-to-service communication
func SetupIntegrationRoutes(router *gin.RouterGroup) {
	integrationServiceURL := getEnv(integrationServiceURLEnvKey, defaultIntegrationServiceURL)

	// Odoo ERP integration
	odoo := router.Group("/odoo")
	{
		// Lead/Customer management
		odoo.POST("/leads", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/leads"))
		odoo.GET("/leads/:id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/leads/:id"))
		odoo.PUT("/leads/:id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/leads/:id"))

		// Convert lead to customer
		odoo.POST("/leads/:id/convert", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/leads/:id/convert"))

		// Customer management
		odoo.GET("/customers/:id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/customers/:id"))
		odoo.PUT("/customers/:id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/customers/:id"))

		// Invoice management
		odoo.POST("/invoices", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/invoices"))
		odoo.GET("/invoices/:id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/invoices/:id"))
		odoo.POST("/invoices/:id/pay", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/invoices/:id/pay"))

		// Commission management in Odoo
		odoo.POST("/commissions", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/commissions"))
		odoo.GET("/commissions/:id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/commissions/:id"))

		// Sync operations
		odoo.POST("/sync/registration/:registration_id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/sync/registration/:registration_id"))
		odoo.POST("/sync/commission/:commission_id", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/sync/commission/:commission_id"))
		odoo.POST("/sync/all", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/sync/all"))

		// Odoo connection status
		odoo.GET("/status", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/status"))
		odoo.POST("/test-connection", proxyToService(integrationServiceURL, "/api/v1/integration/odoo/test-connection"))
	}

	// CIPC (Companies and Intellectual Property Commission) integration
	cipc := router.Group("/cipc")
	{
		cipc.POST("/search", proxyToService(integrationServiceURL, "/api/v1/integration/cipc/search"))
		cipc.POST("/verify", proxyToService(integrationServiceURL, "/api/v1/integration/cipc/verify"))
		cipc.GET("/company/:registration_number", proxyToService(integrationServiceURL, "/api/v1/integration/cipc/company/:registration_number"))
		cipc.GET("/status", proxyToService(integrationServiceURL, "/api/v1/integration/cipc/status"))
	}

	// SARS (South African Revenue Service) integration
	sars := router.Group("/sars")
	{
		sars.POST("/verify-vat", proxyToService(integrationServiceURL, "/api/v1/integration/sars/verify-vat"))
		sars.POST("/verify-tax-number", proxyToService(integrationServiceURL, "/api/v1/integration/sars/verify-tax-number"))
		sars.GET("/status", proxyToService(integrationServiceURL, "/api/v1/integration/sars/status"))
	}

	// Payment gateway integrations
	payments := router.Group("/payments")
	{
		// Stripe
		payments.POST("/stripe/create-payment-intent", proxyToService(integrationServiceURL, "/api/v1/integration/payments/stripe/create-payment-intent"))
		payments.POST("/stripe/webhook", proxyToService(integrationServiceURL, "/api/v1/integration/payments/stripe/webhook"))

		// PayFast (South African payment gateway)
		payments.POST("/payfast/create-payment", proxyToService(integrationServiceURL, "/api/v1/integration/payments/payfast/create-payment"))
		payments.POST("/payfast/webhook", proxyToService(integrationServiceURL, "/api/v1/integration/payments/payfast/webhook"))

		// Payment status
		payments.GET("/status/:payment_id", proxyToService(integrationServiceURL, "/api/v1/integration/payments/status/:payment_id"))
	}

	// Email service integration
	email := router.Group("/email")
	{
		email.POST("/send", proxyToService(integrationServiceURL, "/api/v1/integration/email/send"))
		email.POST("/send-template", proxyToService(integrationServiceURL, "/api/v1/integration/email/send-template"))
		email.GET("/templates", proxyToService(integrationServiceURL, "/api/v1/integration/email/templates"))
	}

	// SMS service integration
	sms := router.Group("/sms")
	{
		sms.POST("/send", proxyToService(integrationServiceURL, "/api/v1/integration/sms/send"))
		sms.GET("/status/:message_id", proxyToService(integrationServiceURL, "/api/v1/integration/sms/status/:message_id"))
	}

	// Webhook management
	webhooks := router.Group("/webhooks")
	{
		webhooks.GET("", proxyToService(integrationServiceURL, "/api/v1/integration/webhooks"))
		webhooks.POST("", proxyToService(integrationServiceURL, "/api/v1/integration/webhooks"))
		webhooks.GET("/:id", proxyToService(integrationServiceURL, "/api/v1/integration/webhooks/:id"))
		webhooks.PUT("/:id", proxyToService(integrationServiceURL, "/api/v1/integration/webhooks/:id"))
		webhooks.DELETE("/:id", proxyToService(integrationServiceURL, "/api/v1/integration/webhooks/:id"))
		webhooks.POST("/:id/test", proxyToService(integrationServiceURL, "/api/v1/integration/webhooks/:id/test"))
	}

	// Integration health and status
	router.GET("/health", proxyToService(integrationServiceURL, "/api/v1/integration/health"))
	router.GET("/status", proxyToService(integrationServiceURL, "/api/v1/integration/status"))
}
