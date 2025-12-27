package handlers

import (
	"net/http"

	"github.com/comply360/integration-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
)

// CIPCHandler handles CIPC-related HTTP requests
type CIPCHandler struct {
	service *services.CIPCService
}

// NewCIPCHandler creates a new CIPC handler
func NewCIPCHandler(service *services.CIPCService) *CIPCHandler {
	return &CIPCHandler{
		service: service,
	}
}

// SetupRoutes sets up the CIPC routes
func (h *CIPCHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/search", h.SearchCompany)
	router.GET("/company/:registration_number", h.GetCompanyDetails)
	router.POST("/validate", h.ValidateCompany)
	router.GET("/status/:registration_number", h.CheckStatus)
}

// SearchCompany handles GET /integrations/cipc/search?company_name=xxx
func (h *CIPCHandler) SearchCompany(c *gin.Context) {
	companyName := c.Query("company_name")
	if companyName == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"company_name query parameter is required",
		))
		return
	}

	results, err := h.service.SearchCompany(companyName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to search CIPC",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": len(results),
	})
}

// GetCompanyDetails handles GET /integrations/cipc/company/:registration_number
func (h *CIPCHandler) GetCompanyDetails(c *gin.Context) {
	registrationNumber := c.Param("registration_number")
	if registrationNumber == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"registration_number is required",
		))
		return
	}

	details, err := h.service.GetCompanyDetails(registrationNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to get company details",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, details)
}

// ValidateCompany handles POST /integrations/cipc/validate
func (h *CIPCHandler) ValidateCompany(c *gin.Context) {
	var req struct {
		RegistrationNumber string `json:"registration_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	result, err := h.service.ValidateCompany(req.RegistrationNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to validate company",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, result)
}

// CheckStatus handles GET /integrations/cipc/status/:registration_number
func (h *CIPCHandler) CheckStatus(c *gin.Context) {
	registrationNumber := c.Param("registration_number")
	if registrationNumber == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"registration_number is required",
		))
		return
	}

	status, err := h.service.CheckStatus(registrationNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to check company status",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"registration_number": registrationNumber,
		"status":             status,
	})
}
