package handlers

import (
	"net/http"
	"strconv"

	"github.com/comply360/commission-service/internal/services"
	"github.com/comply360/shared/errors"
	sharedmiddleware "github.com/comply360/shared/middleware"
	"github.com/comply360/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommissionHandler struct {
	service *services.CommissionService
}

func NewCommissionHandler(service *services.CommissionService) *CommissionHandler {
	return &CommissionHandler{
		service: service,
	}
}

// CreateCommission handles POST /commissions
func (h *CommissionHandler) CreateCommission(c *gin.Context) {
	// Get tenant schema from context
	schema, exists := c.Get("tenant_schema")
	if !exists {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant context not found",
		))
		return
	}

	// Get tenant ID from context
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Parse request body
	var req models.CreateCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIErrorWithDetails(
			errors.ErrInvalidInput,
			"Invalid request body",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	// Parse UUIDs
	registrationID, err := uuid.Parse(req.RegistrationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid registration ID",
		))
		return
	}

	agentID, err := uuid.Parse(req.AgentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid agent ID",
		))
		return
	}

	// Create commission
	commission, err := h.service.CreateCommission(
		schema.(string),
		tenantID,
		registrationID,
		agentID,
		req.RegistrationFee,
		req.CommissionRate,
		req.Currency,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to create commission",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, commission)
}

// GetCommission handles GET /commissions/:id
func (h *CommissionHandler) GetCommission(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Parse commission ID
	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid commission ID",
		))
		return
	}

	// Get commission
	commission, err := h.service.GetCommission(schema.(string), tenantID, commissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(
			errors.ErrNotFound,
			"Commission not found",
		))
		return
	}

	c.JSON(http.StatusOK, commission)
}

// ListCommissions handles GET /commissions
func (h *CommissionHandler) ListCommissions(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Parse query parameters
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	var agentID *uuid.UUID
	if agentIDStr := c.Query("agent_id"); agentIDStr != "" {
		aid, err := uuid.Parse(agentIDStr)
		if err == nil {
			agentID = &aid
		}
	}

	var registrationID *uuid.UUID
	if regIDStr := c.Query("registration_id"); regIDStr != "" {
		rid, err := uuid.Parse(regIDStr)
		if err == nil {
			registrationID = &rid
		}
	}

	// Get commissions
	commissions, total, err := h.service.ListCommissions(schema.(string), tenantID, agentID, registrationID, offset, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to list commissions",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   commissions,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

// GetCommissionSummary handles GET /commissions/summary
func (h *CommissionHandler) GetCommissionSummary(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Parse agent ID (required)
	agentIDStr := c.Query("agent_id")
	if agentIDStr == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"agent_id is required",
		))
		return
	}

	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid agent_id",
		))
		return
	}

	// Get currency (optional, defaults to ZAR)
	currency := c.DefaultQuery("currency", models.CurrencyZAR)

	// Get summary
	summary, err := h.service.GetCommissionSummary(schema.(string), tenantID, agentID, currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to get commission summary",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ApproveCommission handles POST /commissions/:id/approve
func (h *CommissionHandler) ApproveCommission(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Get user ID (approver)
	approvedBy, _ := sharedmiddleware.GetUserID(c)

	// Parse commission ID
	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid commission ID",
		))
		return
	}

	// Approve commission
	if err := h.service.ApproveCommission(schema.(string), tenantID, commissionID, approvedBy); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to approve commission",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commission approved successfully",
	})
}

// MarkCommissionPaid handles POST /commissions/:id/pay
func (h *CommissionHandler) MarkCommissionPaid(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Parse commission ID
	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid commission ID",
		))
		return
	}

	// Parse request body
	var req models.PayCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIErrorWithDetails(
			errors.ErrInvalidInput,
			"Invalid request body",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	// Mark as paid
	if err := h.service.MarkCommissionPaid(schema.(string), tenantID, commissionID, req.PaymentReference); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to mark commission as paid",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commission marked as paid successfully",
	})
}

// CancelCommission handles POST /commissions/:id/cancel
func (h *CommissionHandler) CancelCommission(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantID, _ := sharedmiddleware.GetTenantID(c)

	// Parse commission ID
	commissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid commission ID",
		))
		return
	}

	// Cancel commission
	if err := h.service.CancelCommission(schema.(string), tenantID, commissionID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to cancel commission",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commission cancelled successfully",
	})
}

// SetupRoutes sets up the commission routes
func (h *CommissionHandler) SetupRoutes(r *gin.RouterGroup) {
	r.POST("", h.CreateCommission)
	r.GET("", h.ListCommissions)
	r.GET("/summary", h.GetCommissionSummary)
	r.GET("/:id", h.GetCommission)
	r.POST("/:id/approve", h.ApproveCommission)
	r.POST("/:id/pay", h.MarkCommissionPaid)
	r.POST("/:id/cancel", h.CancelCommission)
}
