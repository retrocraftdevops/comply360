package handlers

import (
	"net/http"
	"strconv"

	"github.com/comply360/registration-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/comply360/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegistrationHandler struct {
	service *services.RegistrationService
}

func NewRegistrationHandler(service *services.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		service: service,
	}
}

// CreateRegistration handles POST /registrations
func (h *RegistrationHandler) CreateRegistration(c *gin.Context) {
	// Get tenant schema from context
	schema, exists := c.Get("tenant_schema")
	if !exists {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant context not found",
		))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}

	// Parse request body
	var registration models.Registration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIErrorWithDetails(
			errors.ErrInvalidInput,
			"Invalid request body",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	// Set tenant ID from context
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))
	registration.TenantID = tenantID

	// Set assigned_to from user ID if not provided
	if registration.AssignedTo == nil {
		userIDStr := userID.(string)
		registration.AssignedTo = &userIDStr
	}

	// Create registration
	if err := h.service.CreateRegistration(schema.(string), &registration); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to create registration",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, registration)
}

// GetRegistration handles GET /registrations/:id
func (h *RegistrationHandler) GetRegistration(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse registration ID
	registrationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid registration ID",
		))
		return
	}

	// Get registration
	registration, err := h.service.GetRegistration(schema.(string), tenantID, registrationID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(
			errors.ErrNotFound,
			"Registration not found",
		))
		return
	}

	c.JSON(http.StatusOK, registration)
}

// ListRegistrations handles GET /registrations
func (h *RegistrationHandler) ListRegistrations(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse query parameters
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	// Get registrations
	registrations, total, err := h.service.ListRegistrations(schema.(string), tenantID, offset, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to list registrations",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   registrations,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

// UpdateRegistration handles PUT /registrations/:id
func (h *RegistrationHandler) UpdateRegistration(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse registration ID
	registrationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid registration ID",
		))
		return
	}

	// Parse request body
	var registration models.Registration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIErrorWithDetails(
			errors.ErrInvalidInput,
			"Invalid request body",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	// Set IDs from path and context
	registration.ID = registrationID
	registration.TenantID = tenantID

	// Update registration
	if err := h.service.UpdateRegistration(schema.(string), &registration); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to update registration",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, registration)
}

// DeleteRegistration handles DELETE /registrations/:id
func (h *RegistrationHandler) DeleteRegistration(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse registration ID
	registrationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid registration ID",
		))
		return
	}

	// Delete registration
	if err := h.service.DeleteRegistration(schema.(string), tenantID, registrationID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to delete registration",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration deleted successfully",
	})
}

// SetupRoutes sets up the registration routes
func (h *RegistrationHandler) SetupRoutes(r *gin.RouterGroup) {
	r.POST("", h.CreateRegistration)
	r.GET("", h.ListRegistrations)
	r.GET("/:id", h.GetRegistration)
	r.PUT("/:id", h.UpdateRegistration)
	r.DELETE("/:id", h.DeleteRegistration)
}
