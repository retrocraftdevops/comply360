package handlers

import (
	"net/http"

	"github.com/comply360/auth-service/internal/services"
	"github.com/comply360/shared/errors"
	"github.com/comply360/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	// Get tenant ID from context
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrTenantNotFound,
			"Tenant ID not found in context",
		))
		return
	}

	user, err := h.authService.Register(tenantID, &req)
	if err != nil {
		c.JSON(http.StatusConflict, errors.NewAPIError(
			errors.ErrAlreadyExists,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	// Get tenant ID from context
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrTenantNotFound,
			"Tenant ID not found in context",
		))
		return
	}

	authResponse, err := h.authService.Login(tenantID, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrInvalidCredentials,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	authResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrInvalidToken,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// VerifyEmail handles email verification
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	// Get tenant ID from context
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrTenantNotFound,
			"Tenant ID not found in context",
		))
		return
	}

	if err := h.authService.VerifyEmail(tenantID, req.Token); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// SetupMFA handles MFA setup
func (h *AuthHandler) SetupMFA(c *gin.Context) {
	var req struct {
		Method string `json:"method" binding:"required,oneof=totp sms email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	// Get tenant ID and user ID from context
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"Tenant ID not found in context",
		))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User ID not found in context",
		))
		return
	}

	qrCodeURL, err := h.authService.SetupMFA(tenantID, userID, req.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(
			errors.ErrInternalServer,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_code_url": qrCodeURL,
		"message":     "Scan the QR code with your authenticator app and verify with a code",
	})
}

// VerifyMFA handles MFA verification
func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	// Get tenant ID and user ID from context
	tenantID, err := getTenantID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"Tenant ID not found in context",
		))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewAPIError(
			errors.ErrUnauthorized,
			"User ID not found in context",
		))
		return
	}

	if err := h.authService.VerifyMFA(tenantID, userID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(
			errors.ErrInvalidInput,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "MFA enabled successfully",
	})
}

// Placeholder handlers for other endpoints
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "OAuth not implemented yet",
	})
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "OAuth not implemented yet",
	})
}

func (h *AuthHandler) DisableMFA(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented yet",
	})
}

// Helper functions
func getTenantID(c *gin.Context) (uuid.UUID, error) {
	// Try X-Tenant-ID header
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr != "" {
		return uuid.Parse(tenantIDStr)
	}

	// Try context
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(uuid.UUID); ok {
			return tid, nil
		}
		if tidStr, ok := tenantID.(string); ok {
			return uuid.Parse(tidStr)
		}
	}

	return uuid.Nil, errors.NewAPIError(errors.ErrTenantNotFound, "Tenant ID not found")
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	// Try X-User-ID header
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr != "" {
		return uuid.Parse(userIDStr)
	}

	// Try context
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			return uid, nil
		}
		if uidStr, ok := userID.(string); ok {
			return uuid.Parse(uidStr)
		}
	}

	return uuid.Nil, errors.NewAPIError(errors.ErrUnauthorized, "User ID not found")
}
