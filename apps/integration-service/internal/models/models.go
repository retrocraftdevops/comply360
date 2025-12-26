package models

import (
	"time"

	"github.com/google/uuid"
)

// Registration represents a Comply360 registration for Odoo sync
type Registration struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	RegistrationType   string     `json:"registration_type"`
	RegistrationNumber string     `json:"registration_number"`
	CompanyName        string     `json:"company_name"`
	ContactPerson      string     `json:"contact_person"`
	Email              string     `json:"email"`
	Phone              string     `json:"phone"`
	Address            string     `json:"address"`
	City               string     `json:"city"`
	Country            string     `json:"country"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	OdooLeadID         *int       `json:"odoo_lead_id,omitempty"`
	OdooPartnerID      *int       `json:"odoo_partner_id,omitempty"`
}

// Commission represents a Comply360 commission for Odoo sync
type Commission struct {
	ID                 string     `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	RegistrationID     uuid.UUID  `json:"registration_id"`
	RegistrationNumber string     `json:"registration_number"`
	RegistrationType   string     `json:"registration_type"`
	AgentID            uuid.UUID  `json:"agent_id"`
	AgentName          string     `json:"agent_name"`
	Amount             float64    `json:"amount"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	OdooCommissionID   *int       `json:"odoo_commission_id,omitempty"`
}

// Invoice represents an invoice to be created in Odoo
type Invoice struct {
	PartnerID  int           `json:"partner_id"`
	Reference  string        `json:"reference"`
	Lines      []InvoiceLine `json:"lines"`
}

// InvoiceLine represents a line item in an invoice
type InvoiceLine struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

// SyncRequest represents a request to sync data to Odoo
type SyncRequest struct {
	Type string      `json:"type"` // registration, commission, invoice
	ID   uuid.UUID   `json:"id"`
	Data interface{} `json:"data"`
}

// SyncResponse represents the response from a sync operation
type SyncResponse struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	OdooID   int         `json:"odoo_id,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// OdooLead represents a CRM lead from Odoo
type OdooLead struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	ContactName string  `json:"contact_name"`
	Email       string  `json:"email_from"`
	Phone       string  `json:"phone"`
	Probability float64 `json:"probability"`
	PartnerID   *int    `json:"partner_id"`
}

// OdooPartner represents a partner (customer) from Odoo
type OdooPartner struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	IsCompany bool   `json:"is_company"`
}

// WebhookPayload represents a webhook payload from external systems
type WebhookPayload struct {
	Event     string                 `json:"event"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ConnectionStatus represents the status of Odoo connection
type ConnectionStatus struct {
	Connected     bool      `json:"connected"`
	Version       string    `json:"version,omitempty"`
	Database      string    `json:"database,omitempty"`
	LastChecked   time.Time `json:"last_checked"`
	Error         string    `json:"error,omitempty"`
}
