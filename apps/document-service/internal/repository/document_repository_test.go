package repository

import (
	"testing"

	"github.com/comply360/shared/models"
	testhelpers "github.com/comply360/shared/testing"
	"github.com/google/uuid"
)

func TestDocumentRepository_Create(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID,
			client_id UUID,
			uploaded_by UUID,
			document_type VARCHAR(100) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_size BIGINT NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			storage_path TEXT NOT NULL,
			storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			verified_at TIMESTAMP,
			verified_by UUID,
			ocr_processed BOOLEAN NOT NULL DEFAULT false,
			ocr_text TEXT,
			ai_verified BOOLEAN NOT NULL DEFAULT false,
			ai_verification_score DECIMAL(5,2),
			ai_verification_notes TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewDocumentRepository(tdb.DB)

	document := &models.Document{
		TenantID:        tdb.TenantID,
		DocumentType:    "id_document",
		FileName:        "test_id.pdf",
		FileSize:        1024 * 512, // 512KB
		MimeType:        "application/pdf",
		StoragePath:     "tenants/test/id_document/test_id.pdf",
		StorageProvider: "s3",
		Status:          "pending",
	}

	// Test: Create document
	err = repo.Create(tdb.Schema, document)
	testhelpers.AssertNoError(t, err, "Failed to create document")
	testhelpers.AssertNotEqual(t, uuid.Nil, document.ID, "Document ID should be set")
}

func TestDocumentRepository_GetByID(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID,
			client_id UUID,
			uploaded_by UUID,
			document_type VARCHAR(100) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_size BIGINT NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			storage_path TEXT NOT NULL,
			storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			verified_at TIMESTAMP,
			verified_by UUID,
			ocr_processed BOOLEAN NOT NULL DEFAULT false,
			ocr_text TEXT,
			ai_verified BOOLEAN NOT NULL DEFAULT false,
			ai_verification_score DECIMAL(5,2),
			ai_verification_notes TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewDocumentRepository(tdb.DB)

	document := &models.Document{
		TenantID:        tdb.TenantID,
		DocumentType:    "proof_of_address",
		FileName:        "utility_bill.pdf",
		FileSize:        2048 * 1024, // 2MB
		MimeType:        "application/pdf",
		StoragePath:     "tenants/test/proof_of_address/utility_bill.pdf",
		StorageProvider: "s3",
		Status:          "pending",
	}

	err = repo.Create(tdb.Schema, document)
	testhelpers.AssertNoError(t, err)

	// Test: Get by ID
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, document.ID)
	testhelpers.AssertNoError(t, err, "Failed to get document by ID")
	testhelpers.AssertEqual(t, document.FileName, retrieved.FileName, "File name mismatch")
	testhelpers.AssertEqual(t, document.FileSize, retrieved.FileSize, "File size mismatch")
	testhelpers.AssertEqual(t, document.Status, retrieved.Status, "Status mismatch")
}

func TestDocumentRepository_UpdateStatus(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID,
			client_id UUID,
			uploaded_by UUID,
			document_type VARCHAR(100) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_size BIGINT NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			storage_path TEXT NOT NULL,
			storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			verified_at TIMESTAMP,
			verified_by UUID,
			ocr_processed BOOLEAN NOT NULL DEFAULT false,
			ocr_text TEXT,
			ai_verified BOOLEAN NOT NULL DEFAULT false,
			ai_verification_score DECIMAL(5,2),
			ai_verification_notes TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewDocumentRepository(tdb.DB)

	document := &models.Document{
		TenantID:        tdb.TenantID,
		DocumentType:    "id_document",
		FileName:        "id_card.pdf",
		FileSize:        500 * 1024,
		MimeType:        "application/pdf",
		StoragePath:     "tenants/test/id_document/id_card.pdf",
		StorageProvider: "s3",
		Status:          "pending",
	}

	err = repo.Create(tdb.Schema, document)
	testhelpers.AssertNoError(t, err)

	// Test: Update status to verified
	verifiedBy := uuid.New()
	err = repo.UpdateStatus(tdb.Schema, tdb.TenantID, document.ID, "verified", &verifiedBy)
	testhelpers.AssertNoError(t, err, "Failed to update document status")

	// Verify update
	retrieved, err := repo.GetByID(tdb.Schema, tdb.TenantID, document.ID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, "verified", retrieved.Status, "Status not updated")
	testhelpers.AssertNotNil(t, retrieved.VerifiedAt, "VerifiedAt should be set")
}

func TestDocumentRepository_ListByRegistration(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID,
			client_id UUID,
			uploaded_by UUID,
			document_type VARCHAR(100) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_size BIGINT NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			storage_path TEXT NOT NULL,
			storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			verified_at TIMESTAMP,
			verified_by UUID,
			ocr_processed BOOLEAN NOT NULL DEFAULT false,
			ocr_text TEXT,
			ai_verified BOOLEAN NOT NULL DEFAULT false,
			ai_verification_score DECIMAL(5,2),
			ai_verification_notes TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewDocumentRepository(tdb.DB)

	registrationID := uuid.New()

	// Create multiple documents for the same registration
	for i := 0; i < 3; i++ {
		document := &models.Document{
			TenantID:        tdb.TenantID,
			RegistrationID:  &registrationID,
			DocumentType:    "id_document",
			FileName:        "document_" + string(rune('A'+i)) + ".pdf",
			FileSize:        1024 * 100,
			MimeType:        "application/pdf",
			StoragePath:     "tenants/test/documents/doc.pdf",
			StorageProvider: "s3",
			Status:          "pending",
		}
		err = repo.Create(tdb.Schema, document)
		testhelpers.AssertNoError(t, err)
	}

	// Test: List documents by registration
	documents, err := repo.ListByRegistration(tdb.Schema, tdb.TenantID, registrationID)
	testhelpers.AssertNoError(t, err, "Failed to list documents by registration")
	testhelpers.AssertEqual(t, 3, len(documents), "Should return 3 documents")
}

func TestDocumentRepository_Delete(t *testing.T) {
	tdb := testhelpers.SetupTestDB(t)
	defer tdb.Cleanup(t)

	_, err := tdb.DB.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL,
			registration_id UUID,
			client_id UUID,
			uploaded_by UUID,
			document_type VARCHAR(100) NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_size BIGINT NOT NULL,
			mime_type VARCHAR(100) NOT NULL,
			storage_path TEXT NOT NULL,
			storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			verified_at TIMESTAMP,
			verified_by UUID,
			ocr_processed BOOLEAN NOT NULL DEFAULT false,
			ocr_text TEXT,
			ai_verified BOOLEAN NOT NULL DEFAULT false,
			ai_verification_score DECIMAL(5,2),
			ai_verification_notes TEXT,
			metadata JSONB DEFAULT '{}'::jsonb,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	testhelpers.AssertNoError(t, err)

	repo := NewDocumentRepository(tdb.DB)

	document := &models.Document{
		TenantID:        tdb.TenantID,
		DocumentType:    "id_document",
		FileName:        "test.pdf",
		FileSize:        1024,
		MimeType:        "application/pdf",
		StoragePath:     "test/path.pdf",
		StorageProvider: "s3",
		Status:          "pending",
	}

	err = repo.Create(tdb.Schema, document)
	testhelpers.AssertNoError(t, err)

	// Test: Delete document
	err = repo.Delete(tdb.Schema, tdb.TenantID, document.ID)
	testhelpers.AssertNoError(t, err, "Failed to delete document")

	// Verify deletion (soft delete)
	_, err = repo.GetByID(tdb.Schema, tdb.TenantID, document.ID)
	testhelpers.AssertError(t, err, "Should not find deleted document")
}
