package handlers

import (
	"net/http"
	"strconv"

	"github.com/comply360/shared/errors"
	"github.com/comply360/shared/models"
	"github.com/comply360/tenant-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TenantHandler struct {
	service *services.TenantService
}

func NewTenantHandler(service *services.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

// CreateTenant creates a new tenant
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req models.CreateTenantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	tenant, err := h.service.CreateTenant(&req)
	if err != nil {
		if err.Error() == "subdomain already exists" {
			c.JSON(http.StatusConflict, errors.NewAPIError(
				errors.ErrTenantAlreadyExists,
				err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternal,
			"Failed to create tenant",
		))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tenant": tenant,
		"message": "Tenant created successfully. Call /provision to set up the tenant environment.",
	})
}

// ProvisionTenant provisions a complete tenant environment
func (h *TenantHandler) ProvisionTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid tenant ID",
		))
		return
	}

	if err := h.service.ProvisionTenant(id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternal,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant provisioned successfully",
	})
}

// GetTenant retrieves a tenant by ID
func (h *TenantHandler) GetTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid tenant ID",
		))
		return
	}

	tenant, err := h.service.GetTenant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(
			errors.ErrTenantNotFound,
			"Tenant not found",
		))
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// ListTenants lists all tenants with pagination
func (h *TenantHandler) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	response, err := h.service.ListTenants(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternal,
			"Failed to list tenants",
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTenant updates a tenant
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid tenant ID",
		))
		return
	}

	var req models.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	tenant, err := h.service.UpdateTenant(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternal,
			"Failed to update tenant",
		))
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// DeleteTenant soft deletes a tenant
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid tenant ID",
		))
		return
	}

	if err := h.service.DeleteTenant(id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternal,
			"Failed to delete tenant",
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant deleted successfully",
	})
}
