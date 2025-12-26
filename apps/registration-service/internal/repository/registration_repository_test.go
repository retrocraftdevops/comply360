package repository

import (
	"testing"

	"github.com/comply360/shared/models"
	testhelpers "github.com/comply360/shared/testing"
	"github.com/google/uuid"
)

func TestRegistrationRepository_Create(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	// Create necessary tables
	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			client_id UUID NOT NULL,
			registration_type VARCHAR(100) NOT NULL,
			company_name VARCHAR(255) NOT NULL,
			registration_number VARCHAR(100),
			jurisdiction VARCHAR(2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			assigned_to VARCHAR(255),
			form_data JSONB DEFAULT '{}'::jsonb,
			submitted_at TIMESTAMP,
			approved_at TIMESTAMP,
			deleted_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewRegistrationRepository(tdb.DB)

	registration := &models.Registration{
		TenantID:         tdb.TenantID,
		ClientID:         uuid.New(),
		RegistrationType: "pty_ltd",
		CompanyName:      "Test Company Pty Ltd",
		Jurisdiction:     "ZA",
		Status:           models.RegistrationStatusDraft,
	}

	// Test: Create registration
	err = repo.Create(tdb.Schema, registration)
	testhelpers.AssertNoError(t, err, "Failed to create registration")
	testhelpers.AssertNotEqual(t, uuid.Nil, registration.ID, "Registration ID should be set")
}

func TestRegistrationRepository_GetByID(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			client_id UUID NOT NULL,
			registration_type VARCHAR(100) NOT NULL,
			company_name VARCHAR(255) NOT NULL,
			registration_number VARCHAR(100),
			jurisdiction VARCHAR(2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			submitted_at TIMESTAMP,
			approved_at TIMESTAMP,
			rejected_at TIMESTAMP,
			rejection_reason TEXT,
			assigned_to VARCHAR(255),
			cipc_reference VARCHAR(255),
			dcip_reference VARCHAR(255),
			odoo_lead_id INTEGER,
			odoo_project_id INTEGER,
			odoo_invoice_id INTEGER,
			form_data JSONB DEFAULT '{}'::jsonb,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewRegistrationRepository(tdb.DB)

	registration := &models.Registration{
		TenantID:         tdb.TenantID,
		ClientID:         uuid.New(),
		RegistrationType: "pty_ltd",
		CompanyName:      "Test Company",
		Jurisdiction:     "ZA",
		Status:           models.RegistrationStatusDraft,
	}

	err = repo.Create(tdb.Schema, registration)
	testhelpers.AssertNoError(t, err)

	// Test: Get by ID
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, registration.ID)
	testhelpers.AssertNoError(t, err, "Failed to get registration by ID")
	testhelpers.AssertEqual(t, registration.CompanyName, retrieved.CompanyName, "Company name mismatch")
	testhelpers.AssertEqual(t, registration.Status, retrieved.Status, "Status mismatch")
}

func TestRegistrationRepository_Update(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			client_id UUID NOT NULL,
			registration_type VARCHAR(100) NOT NULL,
			company_name VARCHAR(255) NOT NULL,
			registration_number VARCHAR(100),
			jurisdiction VARCHAR(2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			submitted_at TIMESTAMP,
			approved_at TIMESTAMP,
			rejected_at TIMESTAMP,
			rejection_reason TEXT,
			assigned_to VARCHAR(255),
			cipc_reference VARCHAR(255),
			dcip_reference VARCHAR(255),
			odoo_lead_id INTEGER,
			odoo_project_id INTEGER,
			odoo_invoice_id INTEGER,
			form_data JSONB DEFAULT '{}'::jsonb,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewRegistrationRepository(tdb.DB)

	registration := &models.Registration{
		TenantID:         tdb.TenantID,
		ClientID:         uuid.New(),
		RegistrationType: "pty_ltd",
		CompanyName:      "Test Company",
		Jurisdiction:     "ZA",
		Status:           models.RegistrationStatusDraft,
	}

	err = repo.Create(tdb.Schema, registration)
	testhelpers.AssertNoError(t, err)

	// Test: Update registration
	registration.Status = models.RegistrationStatusSubmitted
	regNumber := "2024/123456/07"
	registration.RegistrationNumber = &regNumber

	err = repo.Update(tdb.Schema, registration)
	testhelpers.AssertNoError(t, err, "Failed to update registration")

	// Verify update
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, registration.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, models.RegistrationStatusSubmitted, retrieved.Status, "Status not updated")
	testhelpers.AssertNotNil(t, retrieved.RegistrationNumber, "Registration number should be set")
	testhelpers.AssertEqual(t, regNumber, *retrieved.RegistrationNumber, "Registration number mismatch")
}

func TestRegistrationRepository_List(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			client_id UUID NOT NULL,
			registration_type VARCHAR(100) NOT NULL,
			company_name VARCHAR(255) NOT NULL,
			registration_number VARCHAR(100),
			jurisdiction VARCHAR(2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			submitted_at TIMESTAMP,
			approved_at TIMESTAMP,
			rejected_at TIMESTAMP,
			rejection_reason TEXT,
			assigned_to VARCHAR(255),
			cipc_reference VARCHAR(255),
			dcip_reference VARCHAR(255),
			odoo_lead_id INTEGER,
			odoo_project_id INTEGER,
			odoo_invoice_id INTEGER,
			form_data JSONB DEFAULT '{}'::jsonb,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewRegistrationRepository(tdb.DB)

	// Create multiple registrations
	for i := 0; i < 5; i++ {
		registration := &models.Registration{
			TenantID:         tdb.TenantID,
			ClientID:         uuid.New(),
			RegistrationType: "pty_ltd",
			CompanyName:      "Company " + string(rune('A'+i)),
			Jurisdiction:     "ZA",
			Status:           models.RegistrationStatusDraft,
		}
		err = repo.Create(tdb.Schema, registration)
		testhelpers.AssertNoError(t, err)
	}

	// Test: List registrations
	registrations, total, err := repo.List(tdb.Schema, tdb.TenantID, 0, 10, "")
	testhelpers.AssertNoError(t, err, "Failed to list registrations")
	testhelpers.AssertEqual(t, 5, len(registrations), "Should return 5 registrations")
	testhelpers.AssertEqual(t, 5, total, "Total count should be 5")
}

func TestRegistrationRepository_Delete(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS registrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			client_id UUID NOT NULL,
			registration_type VARCHAR(100) NOT NULL,
			company_name VARCHAR(255) NOT NULL,
			registration_number VARCHAR(100),
			jurisdiction VARCHAR(2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'draft',
			assigned_to VARCHAR(255),
			form_data JSONB DEFAULT '{}'::jsonb,
			submitted_at TIMESTAMP,
			approved_at TIMESTAMP,
			deleted_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewRegistrationRepository(tdb.DB)

	registration := &models.Registration{
		TenantID:         tdb.TenantID,
		ClientID:         uuid.New(),
		RegistrationType: "pty_ltd",
		CompanyName:      "Test Company",
		Jurisdiction:     "ZA",
		Status:           models.RegistrationStatusDraft,
	}

	err = repo.Create(tdb.Schema, registration)
	testhelpers.AssertNoError(t, err)

	// Test: Delete registration
	err = repo.Delete(tdb.Schema, tdb.TenantID, registration.ID)
	testhelpers.AssertNoError(t, err, "Failed to delete registration")

	// Verify deletion (soft delete)
	_, err = repo.GetByID(tdb.Schema, tdb.TenantID, registration.ID)
	testhelpers.AssertError(t, err, "Should not find deleted registration")
}
