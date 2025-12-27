package router

import (
	"github.com/gin-gonic/gin"
)

const (
	notificationServiceURLEnvKey = "NOTIFICATION_SERVICE_URL"
	defaultNotificationServiceURL = "http://localhost:8087"
)

// SetupNotificationRoutes configures notification service routes
func SetupNotificationRoutes(router *gin.RouterGroup) {
	notificationServiceURL := getEnv(notificationServiceURLEnvKey, defaultNotificationServiceURL)

	// Email notification routes
	email := router.Group("/email")
	{
		// Custom email
		email.POST("/send", proxyToService(notificationServiceURL, "/api/v1/notifications/email/send"))

		// Registration notifications
		email.POST("/registration/created", proxyToService(notificationServiceURL, "/api/v1/notifications/email/registration/created"))
		email.POST("/registration/submitted", proxyToService(notificationServiceURL, "/api/v1/notifications/email/registration/submitted"))
		email.POST("/registration/approved", proxyToService(notificationServiceURL, "/api/v1/notifications/email/registration/approved"))
		email.POST("/registration/rejected", proxyToService(notificationServiceURL, "/api/v1/notifications/email/registration/rejected"))

		// Document notifications
		email.POST("/document/uploaded", proxyToService(notificationServiceURL, "/api/v1/notifications/email/document/uploaded"))
		email.POST("/document/verified", proxyToService(notificationServiceURL, "/api/v1/notifications/email/document/verified"))

		// Commission notifications
		email.POST("/commission/approved", proxyToService(notificationServiceURL, "/api/v1/notifications/email/commission/approved"))
		email.POST("/commission/paid", proxyToService(notificationServiceURL, "/api/v1/notifications/email/commission/paid"))
	}

	// SMS notification routes (future)
	sms := router.Group("/sms")
	{
		sms.POST("/send", proxyToService(notificationServiceURL, "/api/v1/notifications/sms/send"))
	}
}
