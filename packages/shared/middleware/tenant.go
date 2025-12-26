package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/comply360/shared/errors"
	"github.com/comply360/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// Context keys
	TenantIDKey   = "tenant_id"
	TenantKey     = "tenant"
	SubdomainKey  = "subdomain"
)

// TenantMiddleware extracts tenant context from subdomain or header
func TenantMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID uuid.UUID
		var tenant *models.Tenant

		// Method 1: Extract from subdomain
		host := c.Request.Host
		subdomain := extractSubdomain(host)

		if subdomain != "" && subdomain != "www" && subdomain != "api" {
			// Lookup tenant by subdomain
			t, err := getTenantBySubdomain(db, subdomain)
			if err != nil {
				c.JSON(http.StatusNotFound, errors.NewAPIError(
					errors.ErrTenantNotFound,
					"Tenant not found for subdomain: "+subdomain,
				))
				c.Abort()
				return
			}
			tenant = t
			tenantID = t.ID
		}

		// Method 2: Extract from X-Tenant-ID header (for API calls, admin operations)
		if tenantID == uuid.Nil {
			tenantIDHeader := c.GetHeader("X-Tenant-ID")
			if tenantIDHeader != "" {
				id, err := uuid.Parse(tenantIDHeader)
				if err != nil {
					c.JSON(http.StatusBadRequest, errors.NewAPIError(
						errors.ErrInvalidInput,
						"Invalid tenant ID format",
					))
					c.Abort()
					return
				}

				// Lookup tenant by ID
				t, err := getTenantByID(db, id)
				if err != nil {
					c.JSON(http.StatusNotFound, errors.NewAPIError(
						errors.ErrTenantNotFound,
						"Tenant not found",
					))
					c.Abort()
					return
				}
				tenant = t
				tenantID = t.ID
			}
		}

		// Method 3: Extract from JWT token (will be set by AuthMiddleware)
		// This is handled after AuthMiddleware runs

		if tenantID == uuid.Nil {
			c.JSON(http.StatusBadRequest, errors.NewAPIError(
				errors.ErrTenantNotFound,
				"Tenant context not found",
			))
			c.Abort()
			return
		}

		// Check if tenant is active
		if !tenant.IsActive() {
			c.JSON(http.StatusForbidden, errors.NewAPIError(
				errors.ErrTenantSuspended,
				"Tenant account is suspended or inactive",
			))
			c.Abort()
			return
		}

		// Set tenant context in Gin context
		c.Set(TenantIDKey, tenantID)
		c.Set(TenantKey, tenant)
		c.Set(SubdomainKey, subdomain)
		c.Set("tenant_schema", tenant.TenantSchema())

		// Set PostgreSQL session variable for RLS
		if err := setTenantContext(db, tenantID); err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewAPIError(
				errors.ErrInternal,
				"Failed to set tenant context",
			))
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractSubdomain extracts subdomain from hostname
// Example: agentname.comply360.com -> agentname
func extractSubdomain(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		return parts[0]
	}

	// For localhost development: agentname.localhost -> agentname
	if len(parts) == 2 && parts[1] == "localhost" {
		return parts[0]
	}

	return ""
}

// getTenantBySubdomain fetches tenant from database by subdomain
func getTenantBySubdomain(db *sql.DB, subdomain string) (*models.Tenant, error) {
	var tenant models.Tenant

	query := `
		SELECT id, name, subdomain, domain, status, subscription_tier,
		       company_name, contact_email, contact_phone, country,
		       max_users, created_at, updated_at, deleted_at
		FROM public.tenants
		WHERE subdomain = $1 AND deleted_at IS NULL
	`

	err := db.QueryRow(query, subdomain).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Status,
		&tenant.SubscriptionTier,
		&tenant.CompanyName,
		&tenant.ContactEmail,
		&tenant.ContactPhone,
		&tenant.Country,
		&tenant.MaxUsers,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
		&tenant.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tenant not found")
	}
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// getTenantByID fetches tenant from database by ID
func getTenantByID(db *sql.DB, tenantID uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant

	query := `
		SELECT id, name, subdomain, domain, status, subscription_tier,
		       company_name, contact_email, contact_phone, country,
		       max_users, created_at, updated_at, deleted_at
		FROM public.tenants
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := db.QueryRow(query, tenantID).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Subdomain,
		&tenant.Domain,
		&tenant.Status,
		&tenant.SubscriptionTier,
		&tenant.CompanyName,
		&tenant.ContactEmail,
		&tenant.ContactPhone,
		&tenant.Country,
		&tenant.MaxUsers,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
		&tenant.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tenant not found")
	}
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// setTenantContext sets the PostgreSQL session variable for RLS
func setTenantContext(db *sql.DB, tenantID uuid.UUID) error {
	query := `SELECT set_tenant_context($1)`
	_, err := db.Exec(query, tenantID)
	return err
}

// GetTenantID retrieves tenant ID from Gin context
func GetTenantID(c *gin.Context) (uuid.UUID, error) {
	tenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return uuid.Nil, fmt.Errorf("tenant ID not found in context")
	}

	id, ok := tenantID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid tenant ID type in context")
	}

	return id, nil
}

// GetTenant retrieves tenant from Gin context
func GetTenant(c *gin.Context) (*models.Tenant, error) {
	tenant, exists := c.Get(TenantKey)
	if !exists {
		return nil, fmt.Errorf("tenant not found in context")
	}

	t, ok := tenant.(*models.Tenant)
	if !ok {
		return nil, fmt.Errorf("invalid tenant type in context")
	}

	return t, nil
}
