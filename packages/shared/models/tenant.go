package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Name             string     `json:"name" db:"name"`
	Subdomain        string     `json:"subdomain" db:"subdomain"`
	Domain           *string    `json:"domain,omitempty" db:"domain"`
	Status           string     `json:"status" db:"status"`
	SubscriptionTier string     `json:"subscription_tier" db:"subscription_tier"`
	CompanyName      *string    `json:"company_name,omitempty" db:"company_name"`
	ContactEmail     *string    `json:"contact_email,omitempty" db:"contact_email"`
	ContactPhone     *string    `json:"contact_phone,omitempty" db:"contact_phone"`
	Country          *string    `json:"country,omitempty" db:"country"`
	MaxUsers         int        `json:"max_users" db:"max_users"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// TenantStatus constants
const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
	TenantStatusDeleted   = "deleted"
)

// SubscriptionTier constants
const (
	SubscriptionTierStarter      = "starter"
	SubscriptionTierProfessional = "professional"
	SubscriptionTierEnterprise   = "enterprise"
)

// IsActive checks if tenant is active
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive && t.DeletedAt == nil
}

// TenantSchema returns the schema name for this tenant
func (t *Tenant) TenantSchema() string {
	// Remove hyphens from UUID to create valid PostgreSQL schema name
	return "tenant_" + strings.ReplaceAll(t.ID.String(), "-", "")
}

// TenantListResponse represents a paginated list of tenants
type TenantListResponse struct {
	Tenants    []Tenant `json:"tenants"`
	Total      int      `json:"total"`
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
	TotalPages int      `json:"total_pages"`
}
