package repository

import (
	"testing"

	"github.com/comply360/shared/models"
	testhelpers "github.com/comply360/shared/testing"
	"github.com/google/uuid"
)

func TestTenantRepository_Create(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewTenantRepository(tdb.DB)

	companyName := "Test Company"
	contactEmail := "test@example.com"
	contactPhone := "+27123456789"
	country := "ZA"

	tenant := &models.Tenant{
		Name:             "Test Tenant",
		Subdomain:        "test-" + uuid.New().String()[:8],
		Status:           models.TenantStatusActive,
		SubscriptionTier: models.SubscriptionTierStarter,
		CompanyName:      &companyName,
		ContactEmail:     &contactEmail,
		ContactPhone:     &contactPhone,
		Country:          &country,
		MaxUsers:         10,
	}

	// Test: Create tenant
	err = repo.Create(tenant)
	testhelpers.AssertNoError(t, err, "Failed to create tenant")
	testhelpers.AssertNotEqual(t, uuid.Nil, tenant.ID, "Tenant ID should be set")
}

func TestTenantRepository_GetByID(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewTenantRepository(tdb.DB)

	companyName := "Get By ID Company"
	tenant := &models.Tenant{
		Name:             "Get By ID Tenant",
		Subdomain:        "getbyid-" + uuid.New().String()[:8],
		Status:           models.TenantStatusActive,
		SubscriptionTier: models.SubscriptionTierProfessional,
		CompanyName:      &companyName,
		MaxUsers:         25,
	}

	err = repo.Create(tenant)
	testhelpers.AssertNoError(t, err)

	// Test: Get by ID
	retrieved, err := repo.GetByID(tenant.ID)
	testhelpers.AssertNoError(t, err, "Failed to get tenant by ID")
	testhelpers.AssertEqual(t, tenant.Name, retrieved.Name, "Name mismatch")
	testhelpers.AssertEqual(t, tenant.Subdomain, retrieved.Subdomain, "Subdomain mismatch")
	testhelpers.AssertEqual(t, tenant.Status, retrieved.Status, "Status mismatch")
}

func TestTenantRepository_GetBySubdomain(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewTenantRepository(tdb.DB)

	subdomain := "subdomain-" + uuid.New().String()[:8]
	tenant := &models.Tenant{
		Name:             "Subdomain Test Tenant",
		Subdomain:        subdomain,
		Status:           models.TenantStatusActive,
		SubscriptionTier: models.SubscriptionTierStarter,
		MaxUsers:         10,
	}

	err = repo.Create(tenant)
	testhelpers.AssertNoError(t, err)

	// Test: Get by subdomain
	retrieved, err := repo.GetBySubdomain(subdomain)
	testhelpers.AssertNoError(t, err, "Failed to get tenant by subdomain")
	testhelpers.AssertEqual(t, tenant.ID, retrieved.ID, "ID mismatch")
	testhelpers.AssertEqual(t, tenant.Name, retrieved.Name, "Name mismatch")
}

func TestTenantRepository_List(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	// Clean up any existing tenants from previous tests
	_, err = tdb.DB.Exec(`TRUNCATE TABLE tenants CASCADE`)
	testhelpers.AssertNoError(t, err)

	repo := NewTenantRepository(tdb.DB)

	// Create multiple tenants
	var tenantIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		tenant := &models.Tenant{
			Name:             "List Test Tenant " + string(rune('A'+i)),
			Subdomain:        "list-" + uuid.New().String()[:8],
			Status:           models.TenantStatusActive,
			SubscriptionTier: models.SubscriptionTierStarter,
			MaxUsers:         10,
		}
		err = repo.Create(tenant)
		testhelpers.AssertNoError(t, err)
		tenantIDs = append(tenantIDs, tenant.ID)
	}

	// Test: List tenants
	tenants, total, err := repo.List(1, 100)
	testhelpers.AssertNoError(t, err, "Failed to list tenants")
	testhelpers.AssertTrue(t, len(tenants) >= 3, "Should return at least 3 tenants")
	testhelpers.AssertTrue(t, total >= 3, "Total count should be at least 3")
}

func TestTenantRepository_Update(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewTenantRepository(tdb.DB)

	tenant := &models.Tenant{
		Name:             "Update Test Tenant",
		Subdomain:        "update-" + uuid.New().String()[:8],
		Status:           models.TenantStatusActive,
		SubscriptionTier: models.SubscriptionTierStarter,
		MaxUsers:         10,
	}

	err = repo.Create(tenant)
	testhelpers.AssertNoError(t, err)

	// Update tenant
	newCompanyName := "Updated Company"
	newContactEmail := "updated@example.com"
	tenant.Name = "Updated Tenant Name"
	tenant.CompanyName = &newCompanyName
	tenant.ContactEmail = &newContactEmail
	tenant.Status = models.TenantStatusSuspended

	// Test: Update tenant
	err = repo.Update(tenant)
	testhelpers.AssertNoError(t, err, "Failed to update tenant")

	// Verify update
	retrieved, err := repo.GetByID(tenant.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, "Updated Tenant Name", retrieved.Name, "Name not updated")
	testhelpers.AssertEqual(t, models.TenantStatusSuspended, retrieved.Status, "Status not updated")
	testhelpers.AssertEqual(t, newCompanyName, *retrieved.CompanyName, "Company name not updated")
}

func TestTenantRepository_Delete(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewTenantRepository(tdb.DB)

	tenant := &models.Tenant{
		Name:             "Delete Test Tenant",
		Subdomain:        "delete-" + uuid.New().String()[:8],
		Status:           models.TenantStatusActive,
		SubscriptionTier: models.SubscriptionTierStarter,
		MaxUsers:         10,
	}

	err = repo.Create(tenant)
	testhelpers.AssertNoError(t, err)

	// Test: Delete tenant
	err = repo.Delete(tenant.ID)
	testhelpers.AssertNoError(t, err, "Failed to delete tenant")

	// Verify deletion (soft delete)
	_, err = repo.GetByID(tenant.ID)
	testhelpers.AssertError(t, err, "Should not find deleted tenant")
}

func TestTenantRepository_RecordSchemaCreation(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) NOT NULL UNIQUE,
			domain VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
			company_name VARCHAR(255),
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			country VARCHAR(2),
			max_users INTEGER NOT NULL DEFAULT 10,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	_, err = tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tenant_schemas (
			id SERIAL PRIMARY KEY,
			tenant_id UUID NOT NULL,
			schema_name VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	// Clean up any existing tenant_schemas from previous tests
	_, _ = tdb.DB.Exec(`TRUNCATE TABLE tenant_schemas CASCADE`)

	repo := NewTenantRepository(tdb.DB)

	tenant := &models.Tenant{
		Name:             "Schema Test Tenant",
		Subdomain:        "schema-" + uuid.New().String()[:8],
		Status:           models.TenantStatusActive,
		SubscriptionTier: models.SubscriptionTierStarter,
		MaxUsers:         10,
	}

	err = repo.Create(tenant)
	testhelpers.AssertNoError(t, err)

	// Test: Record schema creation
	schemaName := tenant.TenantSchema()
	err = repo.RecordSchemaCreation(tenant.ID, schemaName)
	testhelpers.AssertNoError(t, err, "Failed to record schema creation")

	// Verify schema was recorded in public schema
	var count int
	err = tdb.DB.QueryRow("SELECT COUNT(*) FROM public.tenant_schemas WHERE tenant_id = $1", tenant.ID).Scan(&count)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, 1, count, "Schema record should exist")
}
