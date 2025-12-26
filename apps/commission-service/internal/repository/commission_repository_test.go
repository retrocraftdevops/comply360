package repository

import (
	"testing"

	"github.com/comply360/shared/models"
	testhelpers "github.com/comply360/shared/testing"
	"github.com/google/uuid"
)

func TestCommissionRepository_Create(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS commissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID NOT NULL,
			agent_id UUID NOT NULL,
			registration_fee DECIMAL(10,2) NOT NULL,
			commission_rate DECIMAL(5,2) NOT NULL,
			commission_amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			approved_at TIMESTAMP,
			approved_by UUID,
			paid_at TIMESTAMP,
			payment_reference VARCHAR(255),
			odoo_commission_id INTEGER,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewCommissionRepository(tdb.DB)

	commission := &models.Commission{
		TenantID:         tdb.TenantID,
		RegistrationID:   uuid.New(),
		AgentID:          uuid.New(),
		RegistrationFee:  5000.00,
		CommissionRate:   15.0,
		CommissionAmount: 750.00, // 15% of 5000
		Currency:         "ZAR",
		Status:           models.CommissionStatusPending,
	}

	// Test: Create commission
	err = repo.Create(tdb.Schema, commission)
	testhelpers.AssertNoError(t, err, "Failed to create commission")
	testhelpers.AssertNotEqual(t, uuid.Nil, commission.ID, "Commission ID should be set")
}

func TestCommissionRepository_GetByID(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS commissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID NOT NULL,
			agent_id UUID NOT NULL,
			registration_fee DECIMAL(10,2) NOT NULL,
			commission_rate DECIMAL(5,2) NOT NULL,
			commission_amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			approved_at TIMESTAMP,
			approved_by UUID,
			paid_at TIMESTAMP,
			payment_reference VARCHAR(255),
			odoo_commission_id INTEGER,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewCommissionRepository(tdb.DB)

	commission := &models.Commission{
		TenantID:         tdb.TenantID,
		RegistrationID:   uuid.New(),
		AgentID:          uuid.New(),
		RegistrationFee:  10000.00,
		CommissionRate:   10.0,
		CommissionAmount: 1000.00,
		Currency:         "ZAR",
		Status:           models.CommissionStatusPending,
	}

	err = repo.Create(tdb.Schema, commission)
	testhelpers.AssertNoError(t, err)

	// Test: Get by ID
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, commission.ID)
	testhelpers.AssertNoError(t, err, "Failed to get commission by ID")
	testhelpers.AssertEqual(t, commission.CommissionAmount, retrieved.CommissionAmount, "Commission amount mismatch")
	testhelpers.AssertEqual(t, commission.CommissionRate, retrieved.CommissionRate, "Commission rate mismatch")
	testhelpers.AssertEqual(t, commission.Status, retrieved.Status, "Status mismatch")
}

func TestCommissionRepository_UpdateStatus(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS commissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID NOT NULL,
			agent_id UUID NOT NULL,
			registration_fee DECIMAL(10,2) NOT NULL,
			commission_rate DECIMAL(5,2) NOT NULL,
			commission_amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			approved_at TIMESTAMP,
			approved_by UUID,
			paid_at TIMESTAMP,
			payment_reference VARCHAR(255),
			odoo_commission_id INTEGER,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewCommissionRepository(tdb.DB)

	commission := &models.Commission{
		TenantID:         tdb.TenantID,
		RegistrationID:   uuid.New(),
		AgentID:          uuid.New(),
		RegistrationFee:  3000.00,
		CommissionRate:   20.0,
		CommissionAmount: 600.00,
		Currency:         "ZAR",
		Status:           models.CommissionStatusPending,
	}

	err = repo.Create(tdb.Schema, commission)
	testhelpers.AssertNoError(t, err)

	// Test: Approve commission
	approvedBy := uuid.New()
	err = repo.Approve(tdb.Schema, tdb.TenantID, commission.ID, approvedBy)
	testhelpers.AssertNoError(t, err, "Failed to approve commission")

	// Verify approval
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, commission.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, models.CommissionStatusApproved, retrieved.Status, "Status not updated to approved")
	testhelpers.AssertNotNil(t, retrieved.ApprovedAt, "ApprovedAt should be set")
}

func TestCommissionRepository_Pay(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS commissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID NOT NULL,
			agent_id UUID NOT NULL,
			registration_fee DECIMAL(10,2) NOT NULL,
			commission_rate DECIMAL(5,2) NOT NULL,
			commission_amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			approved_at TIMESTAMP,
			approved_by UUID,
			paid_at TIMESTAMP,
			payment_reference VARCHAR(255),
			odoo_commission_id INTEGER,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewCommissionRepository(tdb.DB)

	commission := &models.Commission{
		TenantID:         tdb.TenantID,
		RegistrationID:   uuid.New(),
		AgentID:          uuid.New(),
		RegistrationFee:  8000.00,
		CommissionRate:   12.5,
		CommissionAmount: 1000.00,
		Currency:         "ZAR",
		Status:           models.CommissionStatusApproved, // Already approved
	}

	err = repo.Create(tdb.Schema, commission)
	testhelpers.AssertNoError(t, err)

	// Update to approved status manually
	_, err = tdb.DB.Exec(`
		UPDATE commissions
		SET status = $1
		WHERE id = $2
	`, models.CommissionStatusApproved, commission.ID)
	testhelpers.AssertNoError(t, err)

	// Test: Mark as paid
	paymentRef := "PAY-20240126-ABC123"
	err = repo.Pay(tdb.Schema, tdb.TenantID, commission.ID, paymentRef)
	testhelpers.AssertNoError(t, err, "Failed to mark commission as paid")

	// Verify payment
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, commission.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, models.CommissionStatusPaid, retrieved.Status, "Status not updated to paid")
	testhelpers.AssertNotNil(t, retrieved.PaidAt, "PaidAt should be set")
	testhelpers.AssertNotNil(t, retrieved.PaymentReference, "PaymentReference should be set")
	testhelpers.AssertEqual(t, paymentRef, *retrieved.PaymentReference, "Payment reference mismatch")
}

func TestCommissionRepository_GetSummary(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS commissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID NOT NULL,
			agent_id UUID NOT NULL,
			registration_fee DECIMAL(10,2) NOT NULL,
			commission_rate DECIMAL(5,2) NOT NULL,
			commission_amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			approved_at TIMESTAMP,
			approved_by UUID,
			paid_at TIMESTAMP,
			payment_reference VARCHAR(255),
			odoo_commission_id INTEGER,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewCommissionRepository(tdb.DB)

	agentID := uuid.New()

	// Create commissions with different statuses
	statuses := []string{
		models.CommissionStatusPending,
		models.CommissionStatusPending,
		models.CommissionStatusApproved,
		models.CommissionStatusPaid,
	}

	for _, status := range statuses {
		commission := &models.Commission{
			TenantID:         tdb.TenantID,
			RegistrationID:   uuid.New(),
			AgentID:          agentID,
			RegistrationFee:  5000.00,
			CommissionRate:   10.0,
			CommissionAmount: 500.00,
			Currency:         "ZAR",
			Status:           status,
		}
		err = repo.Create(tdb.Schema, commission)
		testhelpers.AssertNoError(t, err)
	}

	// Test: Get summary
	summary, err := repo.GetSummary(tdb.Schema, tdb.TenantID, agentID, "ZAR")
	testhelpers.AssertNoError(t, err, "Failed to get commission summary")
	testhelpers.AssertEqual(t, 2000.00, summary.TotalAmount, "Total amount incorrect")
	testhelpers.AssertEqual(t, 1000.00, summary.PendingAmount, "Pending amount incorrect (2 x 500)")
	testhelpers.AssertEqual(t, 500.00, summary.ApprovedAmount, "Approved amount incorrect")
	testhelpers.AssertEqual(t, 500.00, summary.PaidAmount, "Paid amount incorrect")
}

func TestCommissionRepository_List(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS commissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID NOT NULL,
			agent_id UUID NOT NULL,
			registration_fee DECIMAL(10,2) NOT NULL,
			commission_rate DECIMAL(5,2) NOT NULL,
			commission_amount DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) NOT NULL DEFAULT 'ZAR',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			approved_at TIMESTAMP,
			approved_by UUID,
			paid_at TIMESTAMP,
			payment_reference VARCHAR(255),
			odoo_commission_id INTEGER,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewCommissionRepository(tdb.DB)

	// Create multiple commissions
	for i := 0; i < 5; i++ {
		commission := &models.Commission{
			TenantID:         tdb.TenantID,
			RegistrationID:   uuid.New(),
			AgentID:          uuid.New(),
			RegistrationFee:  5000.00,
			CommissionRate:   10.0,
			CommissionAmount: 500.00,
			Currency:         "ZAR",
			Status:           models.CommissionStatusPending,
		}
		err = repo.Create(tdb.Schema, commission)
		testhelpers.AssertNoError(t, err)
	}

	// Test: List commissions
	commissions, total, err := repo.List(tdb.Schema, tdb.TenantID, nil, nil, 0, 10, "")
	testhelpers.AssertNoError(t, err, "Failed to list commissions")
	testhelpers.AssertEqual(t, 5, len(commissions), "Should return 5 commissions")
	testhelpers.AssertEqual(t, 5, total, "Total count should be 5")
}
