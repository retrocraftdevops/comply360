package handlers

import (
	"net/http"
	"strconv"

	"github.com/comply360/notification-service/internal/models"
	"github.com/comply360/notification-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService *services.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SetupRoutes sets up the notification routes
func (h *NotificationHandler) SetupRoutes(router *gin.RouterGroup) {
	router.POST("", h.CreateNotification)
	router.GET("", h.GetNotifications)
	router.GET("/unread", h.GetUnreadNotifications)
	router.GET("/unread/count", h.GetUnreadCount)
	router.GET("/:id", h.GetNotification)
	router.POST("/:id/read", h.MarkAsRead)
	router.POST("/read-all", h.MarkAllAsRead)
	router.DELETE("/:id", h.DismissNotification)
	router.DELETE("/clear-all", h.ClearAll)

	// Preferences routes
	router.GET("/preferences", h.GetPreferences)
	router.PUT("/preferences", h.UpdatePreferences)

	// Test route
	router.POST("/test", h.SendTestNotification)
}

// CreateNotification handles POST /notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req models.CreateNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	notification, err := h.notificationService.CreateNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to create notification",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetNotifications handles GET /notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// Get user ID from header (forwarded by API Gateway)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	// Get tenant ID from header (optional for system users)
	tenantID := c.GetHeader("X-Tenant-ID")
	// System users (system_admin, global_admin) may have empty tenant_id

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")
	priority := c.Query("priority")
	isReadStr := c.Query("is_read")

	filters := models.NotificationFilters{
		UserID:   userID,
		TenantID: tenantID,
		Page:     page,
		Limit:    limit,
	}

	if category != "" {
		cat := models.NotificationCategory(category)
		filters.Category = &cat
	}

	if priority != "" {
		prio := models.NotificationPriority(priority)
		filters.Priority = &prio
	}

	if isReadStr != "" {
		isRead := isReadStr == "true"
		filters.IsRead = &isRead
	}

	response, err := h.notificationService.GetUserNotifications(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to get notifications",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUnreadNotifications handles GET /notifications/unread
func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant ID is required",
		))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.notificationService.GetUnreadNotifications(c.Request.Context(), userID, tenantID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to get unread notifications",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUnreadCount handles GET /notifications/unread/count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant ID is required",
		))
		return
	}

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to get unread count",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// GetNotification handles GET /notifications/:id
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	notificationID := c.Param("id")

	notification, err := h.notificationService.GetNotification(c.Request.Context(), notificationID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIErrorWithDetails(
			errors.ErrNotFound,
			"Notification not found",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, notification)
}

// MarkAsRead handles POST /notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	notificationID := c.Param("id")

	if err := h.notificationService.MarkAsRead(c.Request.Context(), notificationID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to mark notification as read",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead handles POST /notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant ID is required",
		))
		return
	}

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to mark all notifications as read",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications marked as read",
	})
}

// DismissNotification handles DELETE /notifications/:id
func (h *NotificationHandler) DismissNotification(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	notificationID := c.Param("id")

	if err := h.notificationService.DismissNotification(c.Request.Context(), notificationID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to dismiss notification",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification dismissed",
	})
}

// ClearAll handles DELETE /notifications/clear-all
func (h *NotificationHandler) ClearAll(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant ID is required",
		))
		return
	}

	if err := h.notificationService.ClearAllNotifications(c.Request.Context(), userID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to clear all notifications",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications cleared",
	})
}

// GetPreferences handles GET /notifications/preferences
func (h *NotificationHandler) GetPreferences(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	prefs, err := h.notificationService.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIErrorWithDetails(
			errors.ErrNotFound,
			"Preferences not found",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences handles PUT /notifications/preferences
func (h *NotificationHandler) UpdatePreferences(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	var req models.UpdatePreferencesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	prefs, err := h.notificationService.UpdatePreferences(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to update preferences",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// SendTestNotification handles POST /notifications/test
func (h *NotificationHandler) SendTestNotification(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant ID is required",
		))
		return
	}

	notification, err := h.notificationService.SendTestNotification(c.Request.Context(), userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send test notification",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, notification)
}
