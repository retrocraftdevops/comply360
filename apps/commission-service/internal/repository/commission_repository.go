package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/comply360/shared/models"
	"github.com/google/uuid"
)

type CommissionRepository struct {
	db *sql.DB
}

func NewCommissionRepository(db *sql.DB) *CommissionRepository {
	return &CommissionRepository{db: db}
}

// Create creates a new commission record
func (r *CommissionRepository) Create(schema string, commission *models.Commission) error {
	fmt.Printf("DEBUG REPO: Schema parameter received: '%s' (len=%d)\n", schema, len(schema))

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Set tenant context for RLS - CRITICAL for FK constraint checks!
	tenantIDSetting := fmt.Sprintf("SET LOCAL app.current_tenant_id = '%s'", commission.TenantID.String())
	fmt.Printf("DEBUG REPO: Executing: %s\n", tenantIDSetting)
	_, err = tx.Exec(tenantIDSetting)
	if err != nil {
		fmt.Printf("DEBUG REPO: ERROR setting tenant context: %v\n", err)
		return fmt.Errorf("failed to set tenant context: %w", err)
	}
	fmt.Printf("DEBUG REPO: ✓ Tenant context set\n")

	// Set global admin flag to bypass RLS for FK checks
	fmt.Printf("DEBUG REPO: Executing: SET LOCAL app.is_global_admin = 'true'\n")
	_, err = tx.Exec("SET LOCAL app.is_global_admin = 'true'")
	if err != nil {
		fmt.Printf("DEBUG REPO: ERROR setting admin context: %v\n", err)
		return fmt.Errorf("failed to set admin context: %w", err)
	}
	fmt.Printf("DEBUG REPO: ✓ Admin context set\n")

	// Verify session variables are set
	var tenantIDCheck, adminCheck string
	err = tx.QueryRow("SELECT current_setting('app.current_tenant_id', true), current_setting('app.is_global_admin', true)").Scan(&tenantIDCheck, &adminCheck)
	if err != nil {
		fmt.Printf("DEBUG REPO: ERROR reading session vars: %v\n", err)
	} else {
		fmt.Printf("DEBUG REPO: Session vars verification - tenant_id=%s, is_global_admin=%s\n", tenantIDCheck, adminCheck)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.commissions (
			tenant_id, registration_id, agent_id,
			registration_fee, commission_rate, commission_amount,
			currency, status, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`, schema)

	metadataJSON, err := json.Marshal(commission.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	fmt.Printf("DEBUG REPO: Query length: %d bytes\n", len(query))
	fmt.Printf("DEBUG REPO: Params: tenantID=%s, regID=%s, agentID=%s\n",
		commission.TenantID, commission.RegistrationID, commission.AgentID)

	err = tx.QueryRow(
		query,
		commission.TenantID,
		commission.RegistrationID,
		commission.AgentID,
		commission.RegistrationFee,
		commission.CommissionRate,
		commission.CommissionAmount,
		commission.Currency,
		commission.Status,
		metadataJSON,
	).Scan(&commission.ID, &commission.CreatedAt, &commission.UpdatedAt)

	if err != nil {
		fmt.Printf("DEBUG REPO ERROR: %v\n", err)
		return fmt.Errorf("database insert failed: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("DEBUG REPO SUCCESS: Created commission ID=%s\n", commission.ID)
	return nil
}

// GetByID retrieves a commission by ID
func (r *CommissionRepository) GetByID(schema string, tenantID, commissionID uuid.UUID) (*models.Commission, error) {
	query := fmt.Sprintf(`
		SELECT id, tenant_id, registration_id, agent_id,
			registration_fee, commission_rate, commission_amount,
			currency, status, approved_at, approved_by,
			paid_at, payment_reference, odoo_commission_id,
			created_at, updated_at, metadata
		FROM %s.commissions
		WHERE id = $1 AND tenant_id = $2
	`, schema)

	commission := &models.Commission{}
	var metadataJSON []byte

	err := r.db.QueryRow(query, commissionID, tenantID).Scan(
		&commission.ID,
		&commission.TenantID,
		&commission.RegistrationID,
		&commission.AgentID,
		&commission.RegistrationFee,
		&commission.CommissionRate,
		&commission.CommissionAmount,
		&commission.Currency,
		&commission.Status,
		&commission.ApprovedAt,
		&commission.ApprovedBy,
		&commission.PaidAt,
		&commission.PaymentReference,
		&commission.OdooCommissionID,
		&commission.CreatedAt,
		&commission.UpdatedAt,
		&metadataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("commission not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get commission: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &commission.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return commission, nil
}

// List retrieves commissions with pagination and filters
func (r *CommissionRepository) List(schema string, tenantID uuid.UUID, agentID *uuid.UUID, registrationID *uuid.UUID, offset, limit int, status string) ([]*models.Commission, int, error) {
	// Build query with optional filters
	whereClause := "WHERE tenant_id = $1"
	args := []interface{}{tenantID}
	argCount := 1

	if agentID != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND agent_id = $%d", argCount)
		args = append(args, agentID)
	}

	if registrationID != nil {
		argCount++
		whereClause += fmt.Sprintf(" AND registration_id = $%d", argCount)
		args = append(args, registrationID)
	}

	if status != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s.commissions %s
	`, schema, whereClause)

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count commissions: %w", err)
	}

	// Get commissions
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
		SELECT id, tenant_id, registration_id, agent_id,
			registration_fee, commission_rate, commission_amount,
			currency, status, approved_at, approved_by,
			paid_at, payment_reference, odoo_commission_id,
			created_at, updated_at, metadata
		FROM %s.commissions
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, schema, whereClause, limitArg, offsetArg)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query commissions: %w", err)
	}
	defer rows.Close()

	var commissions []*models.Commission
	for rows.Next() {
		commission := &models.Commission{}
		var metadataJSON []byte

		err := rows.Scan(
			&commission.ID,
			&commission.TenantID,
			&commission.RegistrationID,
			&commission.AgentID,
			&commission.RegistrationFee,
			&commission.CommissionRate,
			&commission.CommissionAmount,
			&commission.Currency,
			&commission.Status,
			&commission.ApprovedAt,
			&commission.ApprovedBy,
			&commission.PaidAt,
			&commission.PaymentReference,
			&commission.OdooCommissionID,
			&commission.CreatedAt,
			&commission.UpdatedAt,
			&metadataJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan commission: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &commission.Metadata); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		commissions = append(commissions, commission)
	}

	return commissions, total, nil
}

// Update updates a commission
func (r *CommissionRepository) Update(schema string, commission *models.Commission) error {
	query := fmt.Sprintf(`
		UPDATE %s.commissions SET
			registration_fee = $1,
			commission_rate = $2,
			commission_amount = $3,
			currency = $4,
			status = $5,
			approved_at = $6,
			approved_by = $7,
			paid_at = $8,
			payment_reference = $9,
			odoo_commission_id = $10,
			metadata = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND tenant_id = $13
	`, schema)

	metadataJSON, err := json.Marshal(commission.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := r.db.Exec(
		query,
		commission.RegistrationFee,
		commission.CommissionRate,
		commission.CommissionAmount,
		commission.Currency,
		commission.Status,
		commission.ApprovedAt,
		commission.ApprovedBy,
		commission.PaidAt,
		commission.PaymentReference,
		commission.OdooCommissionID,
		metadataJSON,
		commission.ID,
		commission.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update commission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commission not found")
	}

	return nil
}

// GetSummary retrieves commission summary statistics for an agent
func (r *CommissionRepository) GetSummary(schema string, tenantID, agentID uuid.UUID, currency string) (*models.CommissionSummary, error) {
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) as total_commissions,
			COALESCE(SUM(commission_amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN commission_amount ELSE 0 END), 0) as pending_amount,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN commission_amount ELSE 0 END), 0) as approved_amount,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN commission_amount ELSE 0 END), 0) as paid_amount
		FROM %s.commissions
		WHERE tenant_id = $1 AND agent_id = $2 AND currency = $3
	`, schema)

	summary := &models.CommissionSummary{Currency: currency}

	err := r.db.QueryRow(query, tenantID, agentID, currency).Scan(
		&summary.TotalCommissions,
		&summary.TotalAmount,
		&summary.PendingAmount,
		&summary.ApprovedAmount,
		&summary.PaidAmount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get commission summary: %w", err)
	}

	return summary, nil
}

// Approve approves a commission
func (r *CommissionRepository) Approve(schema string, tenantID, commissionID, approvedBy uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.commissions SET
			status = 'approved',
			approved_at = CURRENT_TIMESTAMP,
			approved_by = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3
	`, schema)

	result, err := r.db.Exec(query, approvedBy, commissionID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to approve commission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commission not found")
	}

	return nil
}

// Pay marks a commission as paid
func (r *CommissionRepository) Pay(schema string, tenantID, commissionID uuid.UUID, paymentReference string) error {
	query := fmt.Sprintf(`
		UPDATE %s.commissions SET
			status = 'paid',
			paid_at = CURRENT_TIMESTAMP,
			payment_reference = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND status = 'approved'
	`, schema)

	result, err := r.db.Exec(query, paymentReference, commissionID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to mark commission as paid: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commission not found or not approved")
	}

	return nil
}

// GetByRegistrationID retrieves commission by registration ID
func (r *CommissionRepository) GetByRegistrationID(schema string, tenantID, registrationID uuid.UUID) (*models.Commission, error) {
	query := fmt.Sprintf(`
		SELECT id, tenant_id, registration_id, agent_id,
			registration_fee, commission_rate, commission_amount,
			currency, status, approved_at, approved_by,
			paid_at, payment_reference, odoo_commission_id,
			created_at, updated_at, metadata
		FROM %s.commissions
		WHERE registration_id = $1 AND tenant_id = $2
		LIMIT 1
	`, schema)

	commission := &models.Commission{}
	var metadataJSON []byte

	err := r.db.QueryRow(query, registrationID, tenantID).Scan(
		&commission.ID,
		&commission.TenantID,
		&commission.RegistrationID,
		&commission.AgentID,
		&commission.RegistrationFee,
		&commission.CommissionRate,
		&commission.CommissionAmount,
		&commission.Currency,
		&commission.Status,
		&commission.ApprovedAt,
		&commission.ApprovedBy,
		&commission.PaidAt,
		&commission.PaymentReference,
		&commission.OdooCommissionID,
		&commission.CreatedAt,
		&commission.UpdatedAt,
		&metadataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("commission not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get commission: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &commission.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return commission, nil
}
