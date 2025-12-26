package models

import (
	"time"

	"github.com/google/uuid"
)

// Registration represents a company registration
type Registration struct {
	ID                 uuid.UUID              `json:"id" db:"id"`
	TenantID           uuid.UUID              `json:"tenant_id" db:"tenant_id"`
	ClientID           uuid.UUID              `json:"client_id" db:"client_id"`
	RegistrationType   string                 `json:"registration_type" db:"registration_type"`
	CompanyName        string                 `json:"company_name" db:"company_name"`
	RegistrationNumber *string                `json:"registration_number,omitempty" db:"registration_number"`
	Jurisdiction       string                 `json:"jurisdiction" db:"jurisdiction"`
	Status             string                 `json:"status" db:"status"`
	SubmittedAt        *time.Time             `json:"submitted_at,omitempty" db:"submitted_at"`
	ApprovedAt         *time.Time             `json:"approved_at,omitempty" db:"approved_at"`
	RejectedAt         *time.Time             `json:"rejected_at,omitempty" db:"rejected_at"`
	RejectionReason    *string                `json:"rejection_reason,omitempty" db:"rejection_reason"`
	AssignedTo         *string                `json:"assigned_to,omitempty" db:"assigned_to"`
	CIPCReference      *string                `json:"cipc_reference,omitempty" db:"cipc_reference"`
	DCIPReference      *string                `json:"dcip_reference,omitempty" db:"dcip_reference"`
	OdooLeadID         *int                   `json:"odoo_lead_id,omitempty" db:"odoo_lead_id"`
	OdooProjectID      *int                   `json:"odoo_project_id,omitempty" db:"odoo_project_id"`
	OdooInvoiceID      *int                   `json:"odoo_invoice_id,omitempty" db:"odoo_invoice_id"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
	FormData           map[string]interface{} `json:"form_data,omitempty" db:"form_data"`
	Metadata           map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
}

// RegistrationType constants
const (
	RegistrationTypePtyLtd            = "pty_ltd"
	RegistrationTypeCloseCorporation  = "close_corporation"
	RegistrationTypeBusinessName      = "business_name"
	RegistrationTypeVATRegistration   = "vat_registration"
)

// RegistrationStatus constants
const (
	RegistrationStatusDraft      = "draft"
	RegistrationStatusSubmitted  = "submitted"
	RegistrationStatusInReview   = "in_review"
	RegistrationStatusApproved   = "approved"
	RegistrationStatusRejected   = "rejected"
	RegistrationStatusCancelled  = "cancelled"
	RegistrationStatusCompleted  = "completed"
)

// Jurisdiction constants
const (
	JurisdictionSouthAfrica = "ZA"
	JurisdictionZimbabwe    = "ZW"
)

// IsActive checks if registration is active (not deleted)
func (r *Registration) IsActive() bool {
	return r.DeletedAt == nil
}

// IsPending checks if registration is in a pending state
func (r *Registration) IsPending() bool {
	return r.Status == RegistrationStatusDraft ||
	       r.Status == RegistrationStatusSubmitted ||
	       r.Status == RegistrationStatusInReview
}

// IsCompleted checks if registration is completed
func (r *Registration) IsCompleted() bool {
	return r.Status == RegistrationStatusApproved ||
	       r.Status == RegistrationStatusCompleted
}

// CreateRegistrationRequest represents a request to create a new registration
type CreateRegistrationRequest struct {
	ClientID         uuid.UUID              `json:"client_id" binding:"required"`
	RegistrationType string                 `json:"registration_type" binding:"required,oneof=pty_ltd close_corporation business_name vat_registration"`
	CompanyName      string                 `json:"company_name" binding:"required,min=2,max=255"`
	Jurisdiction     string                 `json:"jurisdiction" binding:"required,oneof=ZA ZW"`
	AssignedTo       *string                `json:"assigned_to,omitempty"`
	FormData         map[string]interface{} `json:"form_data,omitempty"`
}

// UpdateRegistrationRequest represents a request to update a registration
type UpdateRegistrationRequest struct {
	RegistrationType   *string                `json:"registration_type,omitempty" binding:"omitempty,oneof=pty_ltd close_corporation business_name vat_registration"`
	CompanyName        *string                `json:"company_name,omitempty" binding:"omitempty,min=2,max=255"`
	RegistrationNumber *string                `json:"registration_number,omitempty"`
	Jurisdiction       *string                `json:"jurisdiction,omitempty" binding:"omitempty,oneof=ZA ZW"`
	Status             *string                `json:"status,omitempty" binding:"omitempty,oneof=draft submitted in_review approved rejected cancelled completed"`
	RejectionReason    *string                `json:"rejection_reason,omitempty"`
	AssignedTo         *string                `json:"assigned_to,omitempty"`
	CIPCReference      *string                `json:"cipc_reference,omitempty"`
	DCIPReference      *string                `json:"dcip_reference,omitempty"`
	FormData           map[string]interface{} `json:"form_data,omitempty"`
}

// RegistrationListResponse represents a paginated list of registrations
type RegistrationListResponse struct {
	Data   []*Registration `json:"data"`
	Total  int             `json:"total"`
	Offset int             `json:"offset"`
	Limit  int             `json:"limit"`
}
