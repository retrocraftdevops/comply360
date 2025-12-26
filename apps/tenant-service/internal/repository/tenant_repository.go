package repository

import (
	"database/sql"
	"fmt"

	"github.com/comply360/shared/models"
	"github.com/google/uuid"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create creates a new tenant record
func (r *TenantRepository) Create(tenant *models.Tenant) error {
	query := `
		INSERT INTO public.tenants (
			name, subdomain, domain, status, subscription_tier,
			company_name, contact_email, contact_phone, country, max_users
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		tenant.Name,
		tenant.Subdomain,
		tenant.Domain,
		tenant.Status,
		tenant.SubscriptionTier,
		tenant.CompanyName,
		tenant.ContactEmail,
		tenant.ContactPhone,
		tenant.Country,
		tenant.MaxUsers,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

// GetByID retrieves a tenant by ID
func (r *TenantRepository) GetByID(id uuid.UUID) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	query := `
		SELECT id, name, subdomain, domain, status, subscription_tier,
		       company_name, contact_email, contact_phone, country,
		       max_users, created_at, updated_at, deleted_at
		FROM public.tenants
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, id).Scan(
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

	return tenant, err
}

// GetBySubdomain retrieves a tenant by subdomain
func (r *TenantRepository) GetBySubdomain(subdomain string) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	query := `
		SELECT id, name, subdomain, domain, status, subscription_tier,
		       company_name, contact_email, contact_phone, country,
		       max_users, created_at, updated_at, deleted_at
		FROM public.tenants
		WHERE subdomain = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, subdomain).Scan(
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

	return tenant, err
}

// List retrieves all tenants with pagination
func (r *TenantRepository) List(page, perPage int) ([]models.Tenant, int, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM public.tenants WHERE deleted_at IS NULL`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get tenants
	query := `
		SELECT id, name, subdomain, domain, status, subscription_tier,
		       company_name, contact_email, contact_phone, country,
		       max_users, created_at, updated_at
		FROM public.tenants
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, 0, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, total, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(tenant *models.Tenant) error {
	query := `
		UPDATE public.tenants
		SET name = $1, company_name = $2, contact_email = $3,
		    contact_phone = $4, status = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(
		query,
		tenant.Name,
		tenant.CompanyName,
		tenant.ContactEmail,
		tenant.ContactPhone,
		tenant.Status,
		tenant.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

// Delete soft deletes a tenant
func (r *TenantRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE public.tenants
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

// RecordSchemaCreation records the tenant schema creation
func (r *TenantRepository) RecordSchemaCreation(tenantID uuid.UUID, schemaName string) error {
	query := `
		INSERT INTO public.tenant_schemas (tenant_id, schema_name, status)
		VALUES ($1, $2, 'active')
	`

	_, err := r.db.Exec(query, tenantID, schemaName)
	return err
}
