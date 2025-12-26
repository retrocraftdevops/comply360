package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/comply360/document-service/internal/repository"
	"github.com/comply360/shared/models"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DocumentService struct {
	repo       *repository.DocumentRepository
	minioClient *minio.Client
	bucketName  string
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
}

func NewDocumentService(repo *repository.DocumentRepository, minioClient *minio.Client, bucketName string, rabbitConn *amqp.Connection) (*DocumentService, error) {
	ch, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange for document events
	err = ch.ExchangeDeclare(
		"comply360.documents", // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &DocumentService{
		repo:       repo,
		minioClient: minioClient,
		bucketName:  bucketName,
		rabbitConn: rabbitConn,
		rabbitCh:   ch,
	}, nil
}

// UploadDocument uploads a document file and creates a database record
func (s *DocumentService) UploadDocument(
	schema string,
	tenantID uuid.UUID,
	registrationID *uuid.UUID,
	clientID *uuid.UUID,
	uploadedBy uuid.UUID,
	documentType string,
	fileName string,
	fileSize int64,
	contentType string,
	fileReader io.Reader,
) (*models.Document, error) {
	// Generate storage path
	storagePath := s.generateStoragePath(tenantID, documentType, fileName)

	// Upload to MinIO
	ctx := context.Background()
	_, err := s.minioClient.PutObject(
		ctx,
		s.bucketName,
		storagePath,
		fileReader,
		fileSize,
		minio.PutObjectOptions{
			ContentType: contentType,
			UserMetadata: map[string]string{
				"tenant-id":     tenantID.String(),
				"document-type": documentType,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create database record
	document := &models.Document{
		ID:              uuid.New(),
		TenantID:        tenantID,
		RegistrationID:  registrationID,
		ClientID:        clientID,
		UploadedBy:      &uploadedBy,
		DocumentType:    documentType,
		FileName:        fileName,
		FileSize:        fileSize,
		MimeType:        contentType,
		StoragePath:     storagePath,
		StorageProvider: models.StorageProviderMinio,
		Status:          models.DocumentStatusPending,
		OCRProcessed:    false,
		AIVerified:      false,
	}

	if err := s.repo.Create(schema, document); err != nil {
		// Attempt to clean up uploaded file if database insert fails
		_ = s.minioClient.RemoveObject(ctx, s.bucketName, storagePath, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	// Publish event
	if err := s.publishEvent("document.uploaded", document); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return document, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(schema string, tenantID, documentID uuid.UUID) (*models.Document, error) {
	return s.repo.GetByID(schema, tenantID, documentID)
}

// GetDocumentDownloadURL generates a presigned URL for downloading a document
func (s *DocumentService) GetDocumentDownloadURL(document *models.Document, expiryDuration time.Duration) (string, error) {
	ctx := context.Background()
	presignedURL, err := s.minioClient.PresignedGetObject(
		ctx,
		s.bucketName,
		document.StoragePath,
		expiryDuration,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return presignedURL.String(), nil
}

// ListDocuments retrieves documents with pagination and filters
func (s *DocumentService) ListDocuments(schema string, tenantID uuid.UUID, registrationID *uuid.UUID, offset, limit int, status, documentType string) ([]*models.Document, int, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(schema, tenantID, registrationID, offset, limit, status, documentType)
}

// UpdateDocument updates a document
func (s *DocumentService) UpdateDocument(schema string, document *models.Document) error {
	// Validate required fields
	if document.ID == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	if document.TenantID == uuid.Nil {
		return fmt.Errorf("tenant_id is required")
	}

	// Update in database
	if err := s.repo.Update(schema, document); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	// Publish event
	if err := s.publishEvent("document.updated", document); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// VerifyDocument marks a document as verified
func (s *DocumentService) VerifyDocument(schema string, tenantID, documentID, verifiedBy uuid.UUID) error {
	document, err := s.repo.GetByID(schema, tenantID, documentID)
	if err != nil {
		return err
	}

	now := time.Now()
	document.Status = models.DocumentStatusVerified
	document.VerifiedAt = &now
	document.VerifiedBy = &verifiedBy

	if err := s.repo.Update(schema, document); err != nil {
		return fmt.Errorf("failed to verify document: %w", err)
	}

	// Publish event
	if err := s.publishEvent("document.verified", document); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// DeleteDocument soft deletes a document (does not remove from storage)
func (s *DocumentService) DeleteDocument(schema string, tenantID, documentID uuid.UUID) error {
	if err := s.repo.Delete(schema, tenantID, documentID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Publish event
	event := map[string]interface{}{
		"document_id": documentID.String(),
		"tenant_id":   tenantID.String(),
	}
	if err := s.publishEvent("document.deleted", event); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// generateStoragePath generates a storage path for a document
func (s *DocumentService) generateStoragePath(tenantID uuid.UUID, documentType, fileName string) string {
	// Generate path: tenants/{tenant_id}/{document_type}/{uuid}_{filename}
	fileExt := filepath.Ext(fileName)
	fileID := uuid.New().String()
	newFileName := fileID + fileExt

	return fmt.Sprintf("tenants/%s/%s/%s", tenantID.String(), documentType, newFileName)
}

// publishEvent publishes an event to RabbitMQ
func (s *DocumentService) publishEvent(eventType string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.rabbitCh.PublishWithContext(
		ctx,
		"comply360.documents", // exchange
		eventType,             // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close closes the service connections
func (s *DocumentService) Close() error {
	if s.rabbitCh != nil {
		return s.rabbitCh.Close()
	}
	return nil
}
