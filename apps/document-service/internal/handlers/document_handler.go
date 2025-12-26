package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/comply360/document-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/comply360/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxUploadSize = 50 << 20 // 50MB

type DocumentHandler struct {
	service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: service,
	}
}

// UploadDocument handles POST /documents
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	// Get tenant schema from context
	schema, exists := c.Get("tenant_schema")
	if !exists {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Tenant context not found",
		))
		return
	}

	// Get tenant ID and user ID from context
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User not authenticated",
		))
		return
	}
	uploadedBy, _ := uuid.Parse(userIDVal.(string))

	// Limit upload size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"File too large or invalid form data",
		))
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"File is required",
		))
		return
	}
	defer file.Close()

	// Get form fields
	documentType := c.PostForm("document_type")
	if documentType == "" {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"document_type is required",
		))
		return
	}

	// Parse optional UUIDs
	var registrationID *uuid.UUID
	if regIDStr := c.PostForm("registration_id"); regIDStr != "" {
		regID, err := uuid.Parse(regIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewAPIError(
				errors.ErrInvalidInput,
				"Invalid registration_id",
			))
			return
		}
		registrationID = &regID
	}

	var clientID *uuid.UUID
	if cIDStr := c.PostForm("client_id"); cIDStr != "" {
		cID, err := uuid.Parse(cIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.NewAPIError(
				errors.ErrInvalidInput,
				"Invalid client_id",
			))
			return
		}
		clientID = &cID
	}

	// Upload document
	document, err := h.service.UploadDocument(
		schema.(string),
		tenantID,
		registrationID,
		clientID,
		uploadedBy,
		documentType,
		header.Filename,
		header.Size,
		header.Header.Get("Content-Type"),
		file,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to upload document",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	// Generate download URL
	downloadURL, _ := h.service.GetDocumentDownloadURL(document, 15*time.Minute)

	c.JSON(http.StatusCreated, models.DocumentUploadResponse{
		Document:    document,
		DownloadURL: downloadURL,
	})
}

// GetDocument handles GET /documents/:id
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse document ID
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid document ID",
		))
		return
	}

	// Get document
	document, err := h.service.GetDocument(schema.(string), tenantID, documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(
			errors.ErrNotFound,
			"Document not found",
		))
		return
	}

	c.JSON(http.StatusOK, document)
}

// GetDocumentDownloadURL handles GET /documents/:id/download
func (h *DocumentHandler) GetDocumentDownloadURL(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse document ID
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid document ID",
		))
		return
	}

	// Get document
	document, err := h.service.GetDocument(schema.(string), tenantID, documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(
			errors.ErrNotFound,
			"Document not found",
		))
		return
	}

	// Generate download URL
	downloadURL, err := h.service.GetDocumentDownloadURL(document, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to generate download URL",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": downloadURL,
		"expires_in":   900, // 15 minutes in seconds
	})
}

// ListDocuments handles GET /documents
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse query parameters
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	documentType := c.Query("document_type")

	var registrationID *uuid.UUID
	if regIDStr := c.Query("registration_id"); regIDStr != "" {
		regID, err := uuid.Parse(regIDStr)
		if err == nil {
			registrationID = &regID
		}
	}

	// Get documents
	documents, total, err := h.service.ListDocuments(schema.(string), tenantID, registrationID, offset, limit, status, documentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to list documents",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   documents,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

// UpdateDocument handles PUT /documents/:id
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse document ID
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid document ID",
		))
		return
	}

	// Get existing document
	document, err := h.service.GetDocument(schema.(string), tenantID, documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(
			errors.ErrNotFound,
			"Document not found",
		))
		return
	}

	// Parse request body
	var req models.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIErrorWithDetails(
			errors.ErrInvalidInput,
			"Invalid request body",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	// Update fields
	if req.DocumentType != nil {
		document.DocumentType = *req.DocumentType
	}
	if req.Status != nil {
		document.Status = *req.Status
	}
	if req.VerifiedBy != nil {
		document.VerifiedBy = req.VerifiedBy
	}
	if req.AIVerificationScore != nil {
		document.AIVerificationScore = req.AIVerificationScore
	}
	if req.AIVerificationNotes != nil {
		document.AIVerificationNotes = req.AIVerificationNotes
	}

	// Update document
	if err := h.service.UpdateDocument(schema.(string), document); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to update document",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, document)
}

// VerifyDocument handles POST /documents/:id/verify
func (h *DocumentHandler) VerifyDocument(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Get user ID
	userIDVal, _ := c.Get("user_id")
	verifiedBy, _ := uuid.Parse(userIDVal.(string))

	// Parse document ID
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid document ID",
		))
		return
	}

	// Verify document
	if err := h.service.VerifyDocument(schema.(string), tenantID, documentID, verifiedBy); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to verify document",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document verified successfully",
	})
}

// DeleteDocument handles DELETE /documents/:id
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	// Get tenant context
	schema, _ := c.Get("tenant_schema")
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := uuid.Parse(tenantIDVal.(string))

	// Parse document ID
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid document ID",
		))
		return
	}

	// Delete document
	if err := h.service.DeleteDocument(schema.(string), tenantID, documentID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to delete document",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}

// SetupRoutes sets up the document routes
func (h *DocumentHandler) SetupRoutes(r *gin.RouterGroup) {
	r.POST("", h.UploadDocument)
	r.GET("", h.ListDocuments)
	r.GET("/:id", h.GetDocument)
	r.GET("/:id/download", h.GetDocumentDownloadURL)
	r.PUT("/:id", h.UpdateDocument)
	r.POST("/:id/verify", h.VerifyDocument)
	r.DELETE("/:id", h.DeleteDocument)
}
