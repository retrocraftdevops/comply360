package handlers

import (
	"net/http"

	"github.com/comply360/notification-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
)

// EmailHandler handles email-related HTTP requests
type EmailHandler struct {
	emailService *services.EmailService
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// SetupRoutes sets up the email routes
func (h *EmailHandler) SetupRoutes(router *gin.RouterGroup) {
	router.POST("/send", h.SendEmail)
	router.POST("/registration/created", h.SendRegistrationCreated)
	router.POST("/registration/submitted", h.SendRegistrationSubmitted)
	router.POST("/registration/approved", h.SendRegistrationApproved)
	router.POST("/registration/rejected", h.SendRegistrationRejected)
	router.POST("/document/uploaded", h.SendDocumentUploaded)
	router.POST("/document/verified", h.SendDocumentVerified)
	router.POST("/commission/approved", h.SendCommissionApproved)
	router.POST("/commission/paid", h.SendCommissionPaid)
}

// SendEmail handles POST /notifications/email/send
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req struct {
		To      []string `json:"to" binding:"required"`
		Subject string   `json:"subject" binding:"required"`
		Body    string   `json:"body" binding:"required"`
		IsHTML  bool     `json:"is_html"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	msg := services.EmailMessage{
		To:      req.To,
		Subject: req.Subject,
		Body:    req.Body,
		IsHTML:  req.IsHTML,
	}

	if err := h.emailService.SendEmail(msg); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
	})
}

// SendRegistrationCreated handles POST /notifications/email/registration/created
func (h *EmailHandler) SendRegistrationCreated(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		ClientName  string `json:"client_name" binding:"required"`
		CompanyName string `json:"company_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendRegistrationCreatedEmail(req.Email, req.ClientName, req.CompanyName); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration created email sent successfully",
	})
}

// SendRegistrationSubmitted handles POST /notifications/email/registration/submitted
func (h *EmailHandler) SendRegistrationSubmitted(c *gin.Context) {
	var req struct {
		Email              string `json:"email" binding:"required,email"`
		ClientName         string `json:"client_name" binding:"required"`
		CompanyName        string `json:"company_name" binding:"required"`
		RegistrationNumber string `json:"registration_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendRegistrationSubmittedEmail(req.Email, req.ClientName, req.CompanyName, req.RegistrationNumber); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration submitted email sent successfully",
	})
}

// SendRegistrationApproved handles POST /notifications/email/registration/approved
func (h *EmailHandler) SendRegistrationApproved(c *gin.Context) {
	var req struct {
		Email              string `json:"email" binding:"required,email"`
		ClientName         string `json:"client_name" binding:"required"`
		CompanyName        string `json:"company_name" binding:"required"`
		RegistrationNumber string `json:"registration_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendRegistrationApprovedEmail(req.Email, req.ClientName, req.CompanyName, req.RegistrationNumber); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration approved email sent successfully",
	})
}

// SendRegistrationRejected handles POST /notifications/email/registration/rejected
func (h *EmailHandler) SendRegistrationRejected(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		ClientName  string `json:"client_name" binding:"required"`
		CompanyName string `json:"company_name" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendRegistrationRejectedEmail(req.Email, req.ClientName, req.CompanyName, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration rejected email sent successfully",
	})
}

// SendDocumentUploaded handles POST /notifications/email/document/uploaded
func (h *EmailHandler) SendDocumentUploaded(c *gin.Context) {
	var req struct {
		Email        string `json:"email" binding:"required,email"`
		ClientName   string `json:"client_name" binding:"required"`
		DocumentType string `json:"document_type" binding:"required"`
		FileName     string `json:"file_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendDocumentUploadedEmail(req.Email, req.ClientName, req.DocumentType, req.FileName); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document uploaded email sent successfully",
	})
}

// SendDocumentVerified handles POST /notifications/email/document/verified
func (h *EmailHandler) SendDocumentVerified(c *gin.Context) {
	var req struct {
		Email        string `json:"email" binding:"required,email"`
		ClientName   string `json:"client_name" binding:"required"`
		DocumentType string `json:"document_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendDocumentVerifiedEmail(req.Email, req.ClientName, req.DocumentType); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document verified email sent successfully",
	})
}

// SendCommissionApproved handles POST /notifications/email/commission/approved
func (h *EmailHandler) SendCommissionApproved(c *gin.Context) {
	var req struct {
		Email     string  `json:"email" binding:"required,email"`
		AgentName string  `json:"agent_name" binding:"required"`
		Amount    float64 `json:"amount" binding:"required"`
		Currency  string  `json:"currency" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendCommissionApprovedEmail(req.Email, req.AgentName, req.Amount, req.Currency); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commission approved email sent successfully",
	})
}

// SendCommissionPaid handles POST /notifications/email/commission/paid
func (h *EmailHandler) SendCommissionPaid(c *gin.Context) {
	var req struct {
		Email            string  `json:"email" binding:"required,email"`
		AgentName        string  `json:"agent_name" binding:"required"`
		Amount           float64 `json:"amount" binding:"required"`
		Currency         string  `json:"currency" binding:"required"`
		PaymentReference string  `json:"payment_reference" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			"Invalid request body",
		))
		return
	}

	if err := h.emailService.SendCommissionPaidEmail(req.Email, req.AgentName, req.Amount, req.Currency, req.PaymentReference); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIErrorWithDetails(
			errors.ErrInternalServer,
			"Failed to send email",
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commission paid email sent successfully",
	})
}
