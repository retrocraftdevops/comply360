package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/comply360/shared/models"
	"github.com/google/uuid"
)

type RegistrationRepository struct {
	db *sql.DB
}

func NewRegistrationRepository(db *sql.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

// Create creates a new registration
func (r *RegistrationRepository) Create(schema string, registration *models.Registration) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.registrations (
			tenant_id, client_id, registration_type, company_name,
			jurisdiction, status, assigned_to, form_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, schema)

	formDataJSON, err := json.Marshal(registration.FormData)
	if err != nil {
		return fmt.Errorf("failed to marshal form data: %w", err)
	}

	return r.db.QueryRow(
		query,
		registration.TenantID,
		registration.ClientID,
		registration.RegistrationType,
		registration.CompanyName,
		registration.Jurisdiction,
		registration.Status,
		registration.AssignedTo,
		formDataJSON,
	).Scan(&registration.ID, &registration.CreatedAt, &registration.UpdatedAt)
}

// GetByID retrieves a registration by ID
func (r *RegistrationRepository) GetByID(schema string, tenantID, registrationID uuid.UUID) (*models.Registration, error) {
	query := fmt.Sprintf(`
		SELECT id, tenant_id, client_id, registration_type, company_name,
			registration_number, jurisdiction, status, submitted_at, approved_at,
			rejected_at, rejection_reason, assigned_to, cipc_reference,
			dcip_reference, odoo_lead_id, odoo_project_id, odoo_invoice_id,
			created_at, updated_at, deleted_at, form_data
		FROM %s.registrations
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, schema)

	registration := &models.Registration{}
	var formDataJSON []byte

	err := r.db.QueryRow(query, registrationID, tenantID).Scan(
		&registration.ID,
		&registration.TenantID,
		&registration.ClientID,
		&registration.RegistrationType,
		&registration.CompanyName,
		&registration.RegistrationNumber,
		&registration.Jurisdiction,
		&registration.Status,
		&registration.SubmittedAt,
		&registration.ApprovedAt,
		&registration.RejectedAt,
		&registration.RejectionReason,
		&registration.AssignedTo,
		&registration.CIPCReference,
		&registration.DCIPReference,
		&registration.OdooLeadID,
		&registration.OdooProjectID,
		&registration.OdooInvoiceID,
		&registration.CreatedAt,
		&registration.UpdatedAt,
		&registration.DeletedAt,
		&formDataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("registration not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get registration: %w", err)
	}

	// Unmarshal form data
	if len(formDataJSON) > 0 {
		if err := json.Unmarshal(formDataJSON, &registration.FormData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal form data: %w", err)
		}
	}

	return registration, nil
}

// List retrieves registrations with pagination
func (r *RegistrationRepository) List(schema string, tenantID uuid.UUID, offset, limit int, status string) ([]*models.Registration, int, error) {
	// Build query with optional status filter
	whereClause := "WHERE tenant_id = $1 AND deleted_at IS NULL"
	args := []interface{}{tenantID}
	argCount := 1

	if status != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s.registrations %s
	`, schema, whereClause)

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count registrations: %w", err)
	}

	// Get registrations
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
		SELECT id, tenant_id, client_id, registration_type, company_name,
			registration_number, jurisdiction, status, submitted_at, approved_at,
			rejected_at, rejection_reason, assigned_to, cipc_reference,
			dcip_reference, odoo_lead_id, odoo_project_id, odoo_invoice_id,
			created_at, updated_at, deleted_at, form_data
		FROM %s.registrations
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, schema, whereClause, limitArg, offsetArg)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query registrations: %w", err)
	}
	defer rows.Close()

	var registrations []*models.Registration
	for rows.Next() {
		registration := &models.Registration{}
		var formDataJSON []byte

		err := rows.Scan(
			&registration.ID,
			&registration.TenantID,
			&registration.ClientID,
			&registration.RegistrationType,
			&registration.CompanyName,
			&registration.RegistrationNumber,
			&registration.Jurisdiction,
			&registration.Status,
			&registration.SubmittedAt,
			&registration.ApprovedAt,
			&registration.RejectedAt,
			&registration.RejectionReason,
			&registration.AssignedTo,
			&registration.CIPCReference,
			&registration.DCIPReference,
			&registration.OdooLeadID,
			&registration.OdooProjectID,
			&registration.OdooInvoiceID,
			&registration.CreatedAt,
			&registration.UpdatedAt,
			&registration.DeletedAt,
			&formDataJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan registration: %w", err)
		}

		// Unmarshal form data
		if len(formDataJSON) > 0 {
			if err := json.Unmarshal(formDataJSON, &registration.FormData); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal form data: %w", err)
			}
		}

		registrations = append(registrations, registration)
	}

	return registrations, total, nil
}

// Update updates a registration
func (r *RegistrationRepository) Update(schema string, registration *models.Registration) error {
	query := fmt.Sprintf(`
		UPDATE %s.registrations SET
			registration_type = $1,
			company_name = $2,
			registration_number = $3,
			jurisdiction = $4,
			status = $5,
			submitted_at = $6,
			approved_at = $7,
			rejected_at = $8,
			rejection_reason = $9,
			assigned_to = $10,
			cipc_reference = $11,
			dcip_reference = $12,
			odoo_lead_id = $13,
			odoo_project_id = $14,
			odoo_invoice_id = $15,
			form_data = $16,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $17 AND tenant_id = $18 AND deleted_at IS NULL
	`, schema)

	formDataJSON, err := json.Marshal(registration.FormData)
	if err != nil {
		return fmt.Errorf("failed to marshal form data: %w", err)
	}

	result, err := r.db.Exec(
		query,
		registration.RegistrationType,
		registration.CompanyName,
		registration.RegistrationNumber,
		registration.Jurisdiction,
		registration.Status,
		registration.SubmittedAt,
		registration.ApprovedAt,
		registration.RejectedAt,
		registration.RejectionReason,
		registration.AssignedTo,
		registration.CIPCReference,
		registration.DCIPReference,
		registration.OdooLeadID,
		registration.OdooProjectID,
		registration.OdooInvoiceID,
		formDataJSON,
		registration.ID,
		registration.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update registration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("registration not found")
	}

	return nil
}

// Delete soft deletes a registration
func (r *RegistrationRepository) Delete(schema string, tenantID, registrationID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.registrations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, schema)

	result, err := r.db.Exec(query, registrationID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete registration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("registration not found")
	}

	return nil
}
