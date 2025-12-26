package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/comply360/shared/models"
	"github.com/comply360/tenant-service/internal/repository"
	"github.com/google/uuid"
)

type TenantService struct {
	repo *repository.TenantRepository
	db   *sql.DB
}

func NewTenantService(repo *repository.TenantRepository, db *sql.DB) *TenantService {
	return &TenantService{
		repo: repo,
		db:   db,
	}
}

// CreateTenant creates a new tenant
func (s *TenantService) CreateTenant(req *models.CreateTenantRequest) (*models.Tenant, error) {
	// Check if subdomain already exists
	existing, _ := s.repo.GetBySubdomain(req.Subdomain)
	if existing != nil {
		return nil, fmt.Errorf("subdomain already exists")
	}

	// Create tenant
	tenant := &models.Tenant{
		Name:             req.Name,
		Subdomain:        req.Subdomain,
		Status:           models.TenantStatusActive,
		SubscriptionTier: req.SubscriptionTier,
		CompanyName:      &req.CompanyName,
		ContactEmail:     &req.ContactEmail,
		ContactPhone:     &req.ContactPhone,
		Country:          &req.Country,
		MaxUsers:         10, // Default
	}

	if err := s.repo.Create(tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// ProvisionTenant provisions a complete tenant environment
func (s *TenantService) ProvisionTenant(tenantID uuid.UUID) error {
	// Get tenant
	tenant, err := s.repo.GetByID(tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Create tenant schema
	schemaName := tenant.TenantSchema()
	if err := s.createTenantSchema(schemaName); err != nil {
		return fmt.Errorf("failed to create tenant schema: %w", err)
	}

	// Run tenant migrations
	if err := s.runTenantMigrations(schemaName); err != nil {
		// Rollback: drop schema
		s.dropTenantSchema(schemaName)
		return fmt.Errorf("failed to run tenant migrations: %w", err)
	}

	// Record schema creation
	if err := s.repo.RecordSchemaCreation(tenantID, schemaName); err != nil {
		return fmt.Errorf("failed to record schema creation: %w", err)
	}

	// TODO: Create initial admin user
	// TODO: Send welcome email
	// TODO: Configure DNS/subdomain

	return nil
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(id uuid.UUID) (*models.Tenant, error) {
	return s.repo.GetByID(id)
}

// ListTenants lists all tenants with pagination
func (s *TenantService) ListTenants(page, perPage int) (*models.TenantListResponse, error) {
	tenants, total, err := s.repo.List(page, perPage)
	if err != nil {
		return nil, err
	}

	totalPages := (total + perPage - 1) / perPage

	return &models.TenantListResponse{
		Tenants:    tenants,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// UpdateTenant updates a tenant
func (s *TenantService) UpdateTenant(id uuid.UUID, req *models.UpdateTenantRequest) (*models.Tenant, error) {
	tenant, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.ContactEmail != nil {
		tenant.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != nil {
		tenant.ContactPhone = req.ContactPhone
	}
	if req.Status != nil {
		tenant.Status = *req.Status
	}
	if req.SubscriptionTier != nil {
		tenant.SubscriptionTier = *req.SubscriptionTier
	}

	if err := s.repo.Update(tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// DeleteTenant soft deletes a tenant
func (s *TenantService) DeleteTenant(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// createTenantSchema creates a new PostgreSQL schema for the tenant
func (s *TenantService) createTenantSchema(schemaName string) error {
	query := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schemaName)
	_, err := s.db.Exec(query)
	return err
}

// dropTenantSchema drops a tenant schema (used for rollback)
func (s *TenantService) dropTenantSchema(schemaName string) error {
	query := fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, schemaName)
	_, err := s.db.Exec(query)
	return err
}

// runTenantMigrations runs the tenant template migration in the tenant schema
func (s *TenantService) runTenantMigrations(schemaName string) error {
	// Set search path to tenant schema
	_, err := s.db.Exec(fmt.Sprintf(`SET search_path TO %s, public`, schemaName))
	if err != nil {
		return err
	}

	// Read and execute tenant template migration
	migrationPath := filepath.Join("database", "migrations", "002_tenant_template", "up.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := s.db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Run RLS policies
	rlsPath := filepath.Join("database", "migrations", "003_rls_policies", "up.sql")
	rlsContent, err := os.ReadFile(rlsPath)
	if err != nil {
		return fmt.Errorf("failed to read RLS migration: %w", err)
	}

	if _, err := s.db.Exec(string(rlsContent)); err != nil {
		return fmt.Errorf("failed to execute RLS migration: %w", err)
	}

	// Reset search path
	_, err = s.db.Exec(`SET search_path TO public`)
	return err
}
