package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/comply360/registration-service/internal/repository"
	"github.com/comply360/shared/models"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RegistrationService struct {
	repo       *repository.RegistrationRepository
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
}

func NewRegistrationService(repo *repository.RegistrationRepository, rabbitConn *amqp.Connection) (*RegistrationService, error) {
	ch, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange for registration events
	err = ch.ExchangeDeclare(
		"comply360.registrations", // name
		"topic",                    // type
		true,                       // durable
		false,                      // auto-deleted
		false,                      // internal
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RegistrationService{
		repo:       repo,
		rabbitConn: rabbitConn,
		rabbitCh:   ch,
	}, nil
}

// CreateRegistration creates a new registration and publishes event
func (s *RegistrationService) CreateRegistration(schema string, registration *models.Registration) error {
	// Set default values
	if registration.ID == uuid.Nil {
		registration.ID = uuid.New()
	}
	if registration.Status == "" {
		registration.Status = "draft"
	}

	// Validate required fields
	if registration.TenantID == uuid.Nil {
		return fmt.Errorf("tenant_id is required")
	}
	if registration.ClientID == uuid.Nil {
		return fmt.Errorf("client_id is required")
	}
	if registration.RegistrationType == "" {
		return fmt.Errorf("registration_type is required")
	}
	if registration.CompanyName == "" {
		return fmt.Errorf("company_name is required")
	}
	if registration.Jurisdiction == "" {
		return fmt.Errorf("jurisdiction is required")
	}

	// Create in database
	if err := s.repo.Create(schema, registration); err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}

	// Publish event
	if err := s.publishEvent("registration.created", registration); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// GetRegistration retrieves a registration by ID
func (s *RegistrationService) GetRegistration(schema string, tenantID, registrationID uuid.UUID) (*models.Registration, error) {
	return s.repo.GetByID(schema, tenantID, registrationID)
}

// ListRegistrations retrieves registrations with pagination
func (s *RegistrationService) ListRegistrations(schema string, tenantID uuid.UUID, offset, limit int, status string) ([]*models.Registration, int, error) {
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

	return s.repo.List(schema, tenantID, offset, limit, status)
}

// UpdateRegistration updates a registration and publishes event
func (s *RegistrationService) UpdateRegistration(schema string, registration *models.Registration) error {
	// Validate required fields
	if registration.ID == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	if registration.TenantID == uuid.Nil {
		return fmt.Errorf("tenant_id is required")
	}

	// Get existing registration to check status transitions
	existing, err := s.repo.GetByID(schema, registration.TenantID, registration.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing registration: %w", err)
	}

	// Validate status transitions
	if err := s.validateStatusTransition(existing.Status, registration.Status); err != nil {
		return err
	}

	// Update timestamps based on status
	now := time.Now()
	switch registration.Status {
	case "submitted":
		if existing.Status != "submitted" {
			registration.SubmittedAt = &now
		}
	case "approved":
		if existing.Status != "approved" {
			registration.ApprovedAt = &now
		}
	case "rejected":
		if existing.Status != "rejected" {
			registration.RejectedAt = &now
		}
	}

	// Update in database
	if err := s.repo.Update(schema, registration); err != nil {
		return fmt.Errorf("failed to update registration: %w", err)
	}

	// Publish event if status changed
	if existing.Status != registration.Status {
		eventType := fmt.Sprintf("registration.status.%s", registration.Status)
		if err := s.publishEvent(eventType, registration); err != nil {
			fmt.Printf("Warning: Failed to publish event: %v\n", err)
		}
	}

	return nil
}

// DeleteRegistration soft deletes a registration
func (s *RegistrationService) DeleteRegistration(schema string, tenantID, registrationID uuid.UUID) error {
	if err := s.repo.Delete(schema, tenantID, registrationID); err != nil {
		return fmt.Errorf("failed to delete registration: %w", err)
	}

	// Publish event
	event := map[string]interface{}{
		"registration_id": registrationID.String(),
		"tenant_id":       tenantID.String(),
	}
	if err := s.publishEvent("registration.deleted", event); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// validateStatusTransition validates if a status transition is allowed
func (s *RegistrationService) validateStatusTransition(from, to string) error {
	validTransitions := map[string][]string{
		"draft":     {"submitted", "cancelled"},
		"submitted": {"approved", "rejected", "cancelled"},
		"approved":  {"completed", "cancelled"},
		"rejected":  {"draft", "cancelled"},
		"completed": {},
		"cancelled": {},
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return fmt.Errorf("invalid current status: %s", from)
	}

	for _, allowed := range allowedStates {
		if allowed == to {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", from, to)
}

// publishEvent publishes an event to RabbitMQ
func (s *RegistrationService) publishEvent(eventType string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.rabbitCh.PublishWithContext(
		ctx,
		"comply360.registrations", // exchange
		eventType,                  // routing key
		false,                      // mandatory
		false,                      // immediate
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
func (s *RegistrationService) Close() error {
	if s.rabbitCh != nil {
		return s.rabbitCh.Close()
	}
	return nil
}
