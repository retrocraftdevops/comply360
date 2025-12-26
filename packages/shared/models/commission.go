package models

import (
	"time"

	"github.com/google/uuid"
)

// Commission represents an agent commission
type Commission struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	TenantID          uuid.UUID              `json:"tenant_id" db:"tenant_id"`
	RegistrationID    uuid.UUID              `json:"registration_id" db:"registration_id"`
	AgentID           uuid.UUID              `json:"agent_id" db:"agent_id"`
	RegistrationFee   float64                `json:"registration_fee" db:"registration_fee"`
	CommissionRate    float64                `json:"commission_rate" db:"commission_rate"`
	CommissionAmount  float64                `json:"commission_amount" db:"commission_amount"`
	Currency          string                 `json:"currency" db:"currency"`
	Status            string                 `json:"status" db:"status"`
	ApprovedAt        *time.Time             `json:"approved_at,omitempty" db:"approved_at"`
	ApprovedBy        *uuid.UUID             `json:"approved_by,omitempty" db:"approved_by"`
	PaidAt            *time.Time             `json:"paid_at,omitempty" db:"paid_at"`
	PaymentReference  *string                `json:"payment_reference,omitempty" db:"payment_reference"`
	OdooCommissionID  *int                   `json:"odoo_commission_id,omitempty" db:"odoo_commission_id"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
	Metadata          map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
}

// CommissionStatus constants
const (
	CommissionStatusPending   = "pending"
	CommissionStatusApproved  = "approved"
	CommissionStatusPaid      = "paid"
	CommissionStatusCancelled = "cancelled"
)

// Currency constants
const (
	CurrencyZAR = "ZAR" // South African Rand
	CurrencyUSD = "USD" // US Dollar
	CurrencyZWL = "ZWL" // Zimbabwean Dollar
)

// IsPending checks if commission is pending
func (c *Commission) IsPending() bool {
	return c.Status == CommissionStatusPending
}

// IsApproved checks if commission is approved
func (c *Commission) IsApproved() bool {
	return c.Status == CommissionStatusApproved
}

// IsPaid checks if commission is paid
func (c *Commission) IsPaid() bool {
	return c.Status == CommissionStatusPaid
}

// CreateCommissionRequest represents a request to create a new commission
type CreateCommissionRequest struct {
	RegistrationID   uuid.UUID `json:"registration_id" binding:"required"`
	AgentID          uuid.UUID `json:"agent_id" binding:"required"`
	RegistrationFee  float64   `json:"registration_fee" binding:"required,gt=0"`
	CommissionRate   float64   `json:"commission_rate" binding:"required,gte=0,lte=100"`
	Currency         string    `json:"currency" binding:"required,oneof=ZAR USD ZWL"`
}

// UpdateCommissionRequest represents a request to update a commission
type UpdateCommissionRequest struct {
	Status           *string    `json:"status,omitempty" binding:"omitempty,oneof=pending approved paid cancelled"`
	PaymentReference *string    `json:"payment_reference,omitempty"`
	OdooCommissionID *int       `json:"odoo_commission_id,omitempty"`
}

// ApproveCommissionRequest represents a request to approve a commission
type ApproveCommissionRequest struct {
	ApprovedBy uuid.UUID `json:"approved_by" binding:"required"`
}

// PayCommissionRequest represents a request to mark a commission as paid
type PayCommissionRequest struct {
	PaymentReference string `json:"payment_reference" binding:"required"`
}

// CommissionListResponse represents a paginated list of commissions
type CommissionListResponse struct {
	Data   []*Commission `json:"data"`
	Total  int           `json:"total"`
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
}

// CommissionSummary represents commission summary statistics
type CommissionSummary struct {
	TotalCommissions  int     `json:"total_commissions"`
	TotalAmount       float64 `json:"total_amount"`
	PendingAmount     float64 `json:"pending_amount"`
	ApprovedAmount    float64 `json:"approved_amount"`
	PaidAmount        float64 `json:"paid_amount"`
	Currency          string  `json:"currency"`
}
