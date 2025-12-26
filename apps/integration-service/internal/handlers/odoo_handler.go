package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/comply360/integration-service/internal/models"
	"github.com/comply360/integration-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OdooHandler struct {
	odooService *services.OdooService
}

func NewOdooHandler(odooService *services.OdooService) *OdooHandler {
	return &OdooHandler{
		odooService: odooService,
	}
}

// CreateLead creates a CRM lead in Odoo from registration data
func (h *OdooHandler) CreateLead(c *gin.Context) {
	var registration models.Registration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	leadID, err := h.odooService.CreateLeadFromRegistration(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to create lead in Odoo: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.SyncResponse{
		Success: true,
		Message: "Lead created successfully in Odoo",
		OdooID:  leadID,
	})
}

// GetLead retrieves a lead from Odoo by ID
func (h *OdooHandler) GetLead(c *gin.Context) {
	leadIDStr := c.Param("id")
	var leadID int
	if _, err := fmt.Sscanf(leadIDStr, "%d", &leadID); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid lead ID",
		))
		return
	}

	// This would require implementing a GetLead method in OdooService
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Get lead not implemented yet",
	})
}

// UpdateLead updates a lead in Odoo
func (h *OdooHandler) UpdateLead(c *gin.Context) {
	leadIDStr := c.Param("id")
	var leadID int
	if _, err := fmt.Sscanf(leadIDStr, "%d", &leadID); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid lead ID",
		))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	if err := h.odooService.UpdateLead(leadID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to update lead: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.SyncResponse{
		Success: true,
		Message: "Lead updated successfully",
		OdooID:  leadID,
	})
}

// ConvertLeadToCustomer converts a lead to a customer
func (h *OdooHandler) ConvertLeadToCustomer(c *gin.Context) {
	leadIDStr := c.Param("id")
	var leadID int
	if _, err := fmt.Sscanf(leadIDStr, "%d", &leadID); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid lead ID",
		))
		return
	}

	partnerID, err := h.odooService.ConvertLeadToCustomer(leadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to convert lead to customer: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.SyncResponse{
		Success: true,
		Message: "Lead converted to customer successfully",
		OdooID:  partnerID,
		Data: map[string]interface{}{
			"lead_id":    leadID,
			"partner_id": partnerID,
		},
	})
}

// CreateInvoice creates an invoice in Odoo
func (h *OdooHandler) CreateInvoice(c *gin.Context) {
	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	invoiceID, err := h.odooService.CreateInvoice(&invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to create invoice: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.SyncResponse{
		Success: true,
		Message: "Invoice created successfully",
		OdooID:  invoiceID,
	})
}

// CreateCommission creates a commission in Odoo
func (h *OdooHandler) CreateCommission(c *gin.Context) {
	var commission models.Commission
	if err := c.ShouldBindJSON(&commission); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	commissionID, err := h.odooService.CreateCommission(&commission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to create commission: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.SyncResponse{
		Success: true,
		Message: "Commission created successfully in Odoo",
		OdooID:  commissionID,
	})
}

// SyncRegistration syncs a registration to Odoo
func (h *OdooHandler) SyncRegistration(c *gin.Context) {
	registrationIDStr := c.Param("registration_id")
	registrationID, err := uuid.Parse(registrationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid registration ID",
		))
		return
	}

	var registration models.Registration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	registration.ID = registrationID

	// Try to find existing lead
	leadID, err := h.odooService.GetLeadByRegistrationID(registrationID)
	if err != nil {
		// Lead doesn't exist, create new one
		leadID, err = h.odooService.CreateLeadFromRegistration(&registration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewAPIError(
				errors.ErrInternalServer,
				"Failed to sync registration: "+err.Error(),
			))
			return
		}

		c.JSON(http.StatusCreated, models.SyncResponse{
			Success: true,
			Message: "Registration synced as new lead",
			OdooID:  leadID,
		})
		return
	}

	// Update existing lead based on status
	if err := h.odooService.SyncRegistrationStatus(registrationID, registration.Status); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to sync registration status: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.SyncResponse{
		Success: true,
		Message: "Registration status synced successfully",
		OdooID:  leadID,
	})
}

// SyncCommission syncs a commission to Odoo
func (h *OdooHandler) SyncCommission(c *gin.Context) {
	commissionIDStr := c.Param("commission_id")
	commissionID, err := uuid.Parse(commissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid commission ID",
		))
		return
	}

	var commission models.Commission
	if err := c.ShouldBindJSON(&commission); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	commission.ID = commissionID.String()

	odooCommissionID, err := h.odooService.CreateCommission(&commission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			"Failed to sync commission: "+err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.SyncResponse{
		Success: true,
		Message: "Commission synced successfully",
		OdooID:  odooCommissionID,
	})
}

// TestConnection tests the Odoo connection
func (h *OdooHandler) TestConnection(c *gin.Context) {
	if err := h.odooService.TestConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ConnectionStatus{
			Connected:   false,
			Error:       err.Error(),
			LastChecked: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ConnectionStatus{
		Connected:   true,
		LastChecked: time.Now(),
	})
}

// GetStatus returns the Odoo connection status
func (h *OdooHandler) GetStatus(c *gin.Context) {
	// Test connection
	err := h.odooService.TestConnection()

	status := models.ConnectionStatus{
		Connected:   err == nil,
		LastChecked: time.Now(),
	}

	if err != nil {
		status.Error = err.Error()
		c.JSON(http.StatusOK, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
