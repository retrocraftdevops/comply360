package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/comply360/shared/models"
	"github.com/google/uuid"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create creates a new document record
func (r *DocumentRepository) Create(schema string, document *models.Document) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.documents (
			tenant_id, registration_id, client_id, uploaded_by,
			document_type, file_name, file_size, mime_type,
			storage_path, storage_provider, status, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`, schema)

	metadataJSON, err := json.Marshal(document.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return r.db.QueryRow(
		query,
		document.TenantID,
		document.RegistrationID,
		document.ClientID,
		document.UploadedBy,
		document.DocumentType,
		document.FileName,
		document.FileSize,
		document.MimeType,
		document.StoragePath,
		document.StorageProvider,
		document.Status,
		metadataJSON,
	).Scan(&document.ID, &document.CreatedAt, &document.UpdatedAt)
}

// GetByID retrieves a document by ID
func (r *DocumentRepository) GetByID(schema string, tenantID, documentID uuid.UUID) (*models.Document, error) {
	query := fmt.Sprintf(`
		SELECT id, tenant_id, registration_id, client_id, uploaded_by,
			document_type, file_name, file_size, mime_type,
			storage_path, storage_provider, status,
			verified_at, verified_by, ocr_processed, ocr_text,
			ai_verified, ai_verification_score, ai_verification_notes,
			created_at, updated_at, deleted_at, metadata
		FROM %s.documents
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, schema)

	document := &models.Document{}
	var metadataJSON []byte

	err := r.db.QueryRow(query, documentID, tenantID).Scan(
		&document.ID,
		&document.TenantID,
		&document.RegistrationID,
		&document.ClientID,
		&document.UploadedBy,
		&document.DocumentType,
		&document.FileName,
		&document.FileSize,
		&document.MimeType,
		&document.StoragePath,
		&document.StorageProvider,
		&document.Status,
		&document.VerifiedAt,
		&document.VerifiedBy,
		&document.OCRProcessed,
		&document.OCRText,
		&document.AIVerified,
		&document.AIVerificationScore,
		&document.AIVerificationNotes,
		&document.CreatedAt,
		&document.UpdatedAt,
		&document.DeletedAt,
		&metadataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &document.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return document, nil
}

// List retrieves documents with pagination and filters
func (r *DocumentRepository) List(schema string, tenantID uuid.UUID, registrationID *uuid.UUID, offset, limit int, status, documentType string) ([]*models.Document, int, error) {
	// Build query with optional filters
	whereClause := "WHERE tenant_id = $1 AND deleted_at IS NULL"
	args := []interface{}{tenantID}
	argCount := 1

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

	if documentType != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, documentType)
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s.documents %s
	`, schema, whereClause)

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Get documents
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
		SELECT id, tenant_id, registration_id, client_id, uploaded_by,
			document_type, file_name, file_size, mime_type,
			storage_path, storage_provider, status,
			verified_at, verified_by, ocr_processed, ocr_text,
			ai_verified, ai_verification_score, ai_verification_notes,
			created_at, updated_at, deleted_at, metadata
		FROM %s.documents
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, schema, whereClause, limitArg, offsetArg)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		document := &models.Document{}
		var metadataJSON []byte

		err := rows.Scan(
			&document.ID,
			&document.TenantID,
			&document.RegistrationID,
			&document.ClientID,
			&document.UploadedBy,
			&document.DocumentType,
			&document.FileName,
			&document.FileSize,
			&document.MimeType,
			&document.StoragePath,
			&document.StorageProvider,
			&document.Status,
			&document.VerifiedAt,
			&document.VerifiedBy,
			&document.OCRProcessed,
			&document.OCRText,
			&document.AIVerified,
			&document.AIVerificationScore,
			&document.AIVerificationNotes,
			&document.CreatedAt,
			&document.UpdatedAt,
			&document.DeletedAt,
			&metadataJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan document: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &document.Metadata); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		documents = append(documents, document)
	}

	return documents, total, nil
}

// Update updates a document
func (r *DocumentRepository) Update(schema string, document *models.Document) error {
	query := fmt.Sprintf(`
		UPDATE %s.documents SET
			document_type = $1,
			file_name = $2,
			status = $3,
			verified_at = $4,
			verified_by = $5,
			ocr_processed = $6,
			ocr_text = $7,
			ai_verified = $8,
			ai_verification_score = $9,
			ai_verification_notes = $10,
			metadata = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND tenant_id = $13 AND deleted_at IS NULL
	`, schema)

	metadataJSON, err := json.Marshal(document.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := r.db.Exec(
		query,
		document.DocumentType,
		document.FileName,
		document.Status,
		document.VerifiedAt,
		document.VerifiedBy,
		document.OCRProcessed,
		document.OCRText,
		document.AIVerified,
		document.AIVerificationScore,
		document.AIVerificationNotes,
		metadataJSON,
		document.ID,
		document.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}

// Delete soft deletes a document
func (r *DocumentRepository) Delete(schema string, tenantID, documentID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.documents
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, schema)

	result, err := r.db.Exec(query, documentID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}
