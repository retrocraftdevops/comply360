package models

import (
	"time"

	"github.com/google/uuid"
)

// Document represents a document in the system
type Document struct {
	ID                   uuid.UUID              `json:"id" db:"id"`
	TenantID             uuid.UUID              `json:"tenant_id" db:"tenant_id"`
	RegistrationID       *uuid.UUID             `json:"registration_id,omitempty" db:"registration_id"`
	ClientID             *uuid.UUID             `json:"client_id,omitempty" db:"client_id"`
	UploadedBy           *uuid.UUID             `json:"uploaded_by,omitempty" db:"uploaded_by"`
	DocumentType         string                 `json:"document_type" db:"document_type"`
	FileName             string                 `json:"file_name" db:"file_name"`
	FileSize             int64                  `json:"file_size" db:"file_size"`
	MimeType             string                 `json:"mime_type" db:"mime_type"`
	StoragePath          string                 `json:"storage_path" db:"storage_path"`
	StorageProvider      string                 `json:"storage_provider" db:"storage_provider"`
	Status               string                 `json:"status" db:"status"`
	VerifiedAt           *time.Time             `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedBy           *uuid.UUID             `json:"verified_by,omitempty" db:"verified_by"`
	OCRProcessed         bool                   `json:"ocr_processed" db:"ocr_processed"`
	OCRText              *string                `json:"ocr_text,omitempty" db:"ocr_text"`
	AIVerified           bool                   `json:"ai_verified" db:"ai_verified"`
	AIVerificationScore  *float64               `json:"ai_verification_score,omitempty" db:"ai_verification_score"`
	AIVerificationNotes  *string                `json:"ai_verification_notes,omitempty" db:"ai_verification_notes"`
	CreatedAt            time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
	Metadata             map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
}

// DocumentStatus constants
const (
	DocumentStatusPending  = "pending"
	DocumentStatusVerified = "verified"
	DocumentStatusRejected = "rejected"
	DocumentStatusExpired  = "expired"
)

// StorageProvider constants
const (
	StorageProviderS3    = "s3"
	StorageProviderMinio = "minio"
)

// Common document types
const (
	DocumentTypeID               = "id_document"
	DocumentTypeProofOfAddress   = "proof_of_address"
	DocumentTypeCompanyDocs      = "company_documents"
	DocumentTypeBankStatement    = "bank_statement"
	DocumentTypeTaxClearance     = "tax_clearance"
	DocumentTypeCIPCDocuments    = "cipc_documents"
)

// IsActive checks if document is active (not deleted)
func (d *Document) IsActive() bool {
	return d.DeletedAt == nil
}

// IsVerified checks if document is verified
func (d *Document) IsVerified() bool {
	return d.Status == DocumentStatusVerified
}

// UploadDocumentRequest represents a request to upload a document
type UploadDocumentRequest struct {
	RegistrationID *uuid.UUID `form:"registration_id"`
	ClientID       *uuid.UUID `form:"client_id"`
	DocumentType   string     `form:"document_type" binding:"required"`
}

// UpdateDocumentRequest represents a request to update document metadata
type UpdateDocumentRequest struct {
	DocumentType        *string  `json:"document_type,omitempty"`
	Status              *string  `json:"status,omitempty" binding:"omitempty,oneof=pending verified rejected expired"`
	VerifiedBy          *uuid.UUID `json:"verified_by,omitempty"`
	AIVerificationScore *float64 `json:"ai_verification_score,omitempty"`
	AIVerificationNotes *string  `json:"ai_verification_notes,omitempty"`
}

// DocumentListResponse represents a paginated list of documents
type DocumentListResponse struct {
	Data   []*Document `json:"data"`
	Total  int         `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
}

// DocumentUploadResponse represents the response after uploading a document
type DocumentUploadResponse struct {
	Document   *Document `json:"document"`
	UploadURL  string    `json:"upload_url,omitempty"`
	DownloadURL string   `json:"download_url,omitempty"`
}
