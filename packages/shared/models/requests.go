package models

// Auth Requests
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,strong_password"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,strong_password"`
}

// Registration Requests
type CreateRegistrationRequest struct {
	ClientID         string                 `json:"client_id" validate:"required,uuid"`
	RegistrationType string                 `json:"registration_type" validate:"required,registration_type"`
	CompanyName      string                 `json:"company_name" validate:"required,min=2,max=255"`
	Jurisdiction     string                 `json:"jurisdiction" validate:"required,jurisdiction"`
	FormData         map[string]interface{} `json:"form_data,omitempty"`
}

type UpdateRegistrationRequest struct {
	RegistrationNumber *string                `json:"registration_number,omitempty" validate:"omitempty,company_registration_number"`
	Status             *string                `json:"status,omitempty" validate:"omitempty,registration_status"`
	RejectionReason    *string                `json:"rejection_reason,omitempty"`
	FormData           map[string]interface{} `json:"form_data,omitempty"`
}

// Document Requests
type DocumentUploadRequest struct {
	DocumentType   string `form:"document_type" validate:"required,document_type"`
	RegistrationID string `form:"registration_id,omitempty" validate:"omitempty,uuid"`
	ClientID       string `form:"client_id,omitempty" validate:"omitempty,uuid"`
}

// Commission Requests
type CreateCommissionRequest struct {
	RegistrationID  string  `json:"registration_id" validate:"required,uuid"`
	AgentID         string  `json:"agent_id" validate:"required,uuid"`
	RegistrationFee float64 `json:"registration_fee" validate:"required,gt=0"`
	CommissionRate  float64 `json:"commission_rate" validate:"required,commission_rate"`
	Currency        string  `json:"currency" validate:"required,currency"`
}

type PayCommissionRequest struct {
	PaymentReference string `json:"payment_reference" validate:"required,min=5,max=255"`
}

// Client Requests
type CreateClientRequest struct {
	ClientType   string  `json:"client_type" validate:"required,oneof=individual company"`
	FullName     *string `json:"full_name,omitempty" validate:"required_if=ClientType individual,omitempty,min=2,max=255"`
	CompanyName  *string `json:"company_name,omitempty" validate:"required_if=ClientType company,omitempty,min=2,max=255"`
	IDNumber     *string `json:"id_number,omitempty" validate:"omitempty,sa_id_number"`
	TaxNumber    *string `json:"tax_number,omitempty"`
	VATNumber    *string `json:"vat_number,omitempty" validate:"omitempty,vat_number"`
	Email        string  `json:"email" validate:"required,email"`
	Phone        *string `json:"phone,omitempty" validate:"omitempty,phone"`
	Mobile       *string `json:"mobile,omitempty" validate:"omitempty,phone"`
	CountryCode  *string `json:"country_code,omitempty" validate:"omitempty,country_code,len=2"`
}

type UpdateClientRequest struct {
	FullName    *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	CompanyName *string `json:"company_name,omitempty" validate:"omitempty,min=2,max=255"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,phone"`
	Mobile      *string `json:"mobile,omitempty" validate:"omitempty,phone"`
	Status      *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
}

// Tenant Requests
type CreateTenantRequest struct {
	Name              string `json:"name" validate:"required,min=2,max=255"`
	Subdomain         string `json:"subdomain" validate:"required,subdomain,min=3,max=63"`
	CompanyName       string `json:"company_name" validate:"required,min=2,max=255"`
	ContactEmail      string `json:"contact_email" validate:"required,email"`
	ContactPhone      string `json:"contact_phone,omitempty" validate:"omitempty,phone"`
	Country           string `json:"country" validate:"required,country_code,len=2"`
	SubscriptionTier  string `json:"subscription_tier,omitempty" validate:"omitempty,oneof=starter professional enterprise"`
}

type UpdateTenantRequest struct {
	Name             *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	ContactEmail     *string `json:"contact_email,omitempty" validate:"omitempty,email"`
	ContactPhone     *string `json:"contact_phone,omitempty" validate:"omitempty,phone"`
	Status           *string `json:"status,omitempty" validate:"omitempty,oneof=active suspended deleted"`
	SubscriptionTier *string `json:"subscription_tier,omitempty" validate:"omitempty,oneof=starter professional enterprise"`
}

// Pagination Request
type PaginationRequest struct {
	Page  int `form:"page" validate:"omitempty,min=1"`
	Limit int `form:"limit" validate:"omitempty,min=1,max=100"`
}

// Default pagination values
func (p *PaginationRequest) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}

// Calculate offset for database queries
func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
