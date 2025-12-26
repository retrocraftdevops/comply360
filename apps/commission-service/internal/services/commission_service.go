package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/comply360/commission-service/internal/repository"
	"github.com/comply360/shared/models"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CommissionService struct {
	repo       *repository.CommissionRepository
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
}

func NewCommissionService(repo *repository.CommissionRepository, rabbitConn *amqp.Connection) (*CommissionService, error) {
	ch, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange for commission events
	err = ch.ExchangeDeclare(
		"comply360.commissions", // name
		"topic",                  // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &CommissionService{
		repo:       repo,
		rabbitConn: rabbitConn,
		rabbitCh:   ch,
	}, nil
}

// CreateCommission creates a new commission with automatic calculation
func (s *CommissionService) CreateCommission(schema string, tenantID, registrationID, agentID uuid.UUID, registrationFee, commissionRate float64, currency string) (*models.Commission, error) {
	// Validate inputs
	if registrationFee <= 0 {
		return nil, fmt.Errorf("registration fee must be greater than 0")
	}
	if commissionRate < 0 || commissionRate > 100 {
		return nil, fmt.Errorf("commission rate must be between 0 and 100")
	}

	// Calculate commission amount
	commissionAmount := s.calculateCommissionAmount(registrationFee, commissionRate)

	// Create commission
	commission := &models.Commission{
		ID:               uuid.New(),
		TenantID:         tenantID,
		RegistrationID:   registrationID,
		AgentID:          agentID,
		RegistrationFee:  registrationFee,
		CommissionRate:   commissionRate,
		CommissionAmount: commissionAmount,
		Currency:         currency,
		Status:           models.CommissionStatusPending,
	}

	if err := s.repo.Create(schema, commission); err != nil {
		return nil, fmt.Errorf("failed to create commission: %w", err)
	}

	// Publish event
	if err := s.publishEvent("commission.created", commission); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return commission, nil
}

// GetCommission retrieves a commission by ID
func (s *CommissionService) GetCommission(schema string, tenantID, commissionID uuid.UUID) (*models.Commission, error) {
	return s.repo.GetByID(schema, tenantID, commissionID)
}

// GetCommissionByRegistration retrieves commission for a specific registration
func (s *CommissionService) GetCommissionByRegistration(schema string, tenantID, registrationID uuid.UUID) (*models.Commission, error) {
	return s.repo.GetByRegistrationID(schema, tenantID, registrationID)
}

// ListCommissions retrieves commissions with pagination and filters
func (s *CommissionService) ListCommissions(schema string, tenantID uuid.UUID, agentID *uuid.UUID, registrationID *uuid.UUID, offset, limit int, status string) ([]*models.Commission, int, error) {
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

	return s.repo.List(schema, tenantID, agentID, registrationID, offset, limit, status)
}

// GetCommissionSummary retrieves commission summary for an agent
func (s *CommissionService) GetCommissionSummary(schema string, tenantID, agentID uuid.UUID, currency string) (*models.CommissionSummary, error) {
	if currency == "" {
		currency = models.CurrencyZAR
	}

	return s.repo.GetSummary(schema, tenantID, agentID, currency)
}

// ApproveCommission approves a commission
func (s *CommissionService) ApproveCommission(schema string, tenantID, commissionID, approvedBy uuid.UUID) error {
	commission, err := s.repo.GetByID(schema, tenantID, commissionID)
	if err != nil {
		return err
	}

	if commission.Status != models.CommissionStatusPending {
		return fmt.Errorf("commission must be in pending status to approve")
	}

	now := time.Now()
	commission.Status = models.CommissionStatusApproved
	commission.ApprovedAt = &now
	commission.ApprovedBy = &approvedBy

	if err := s.repo.Update(schema, commission); err != nil {
		return fmt.Errorf("failed to approve commission: %w", err)
	}

	// Publish event
	if err := s.publishEvent("commission.approved", commission); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// MarkCommissionPaid marks a commission as paid
func (s *CommissionService) MarkCommissionPaid(schema string, tenantID, commissionID uuid.UUID, paymentReference string) error {
	commission, err := s.repo.GetByID(schema, tenantID, commissionID)
	if err != nil {
		return err
	}

	if commission.Status != models.CommissionStatusApproved {
		return fmt.Errorf("commission must be approved before marking as paid")
	}

	now := time.Now()
	commission.Status = models.CommissionStatusPaid
	commission.PaidAt = &now
	commission.PaymentReference = &paymentReference

	if err := s.repo.Update(schema, commission); err != nil {
		return fmt.Errorf("failed to mark commission as paid: %w", err)
	}

	// Publish event
	if err := s.publishEvent("commission.paid", commission); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// CancelCommission cancels a commission
func (s *CommissionService) CancelCommission(schema string, tenantID, commissionID uuid.UUID) error {
	commission, err := s.repo.GetByID(schema, tenantID, commissionID)
	if err != nil {
		return err
	}

	if commission.Status == models.CommissionStatusPaid {
		return fmt.Errorf("cannot cancel a paid commission")
	}

	commission.Status = models.CommissionStatusCancelled

	if err := s.repo.Update(schema, commission); err != nil {
		return fmt.Errorf("failed to cancel commission: %w", err)
	}

	// Publish event
	if err := s.publishEvent("commission.cancelled", commission); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// RecalculateCommission recalculates commission amount based on new fee or rate
func (s *CommissionService) RecalculateCommission(schema string, tenantID, commissionID uuid.UUID, newRegistrationFee *float64, newCommissionRate *float64) error {
	commission, err := s.repo.GetByID(schema, tenantID, commissionID)
	if err != nil {
		return err
	}

	if commission.Status != models.CommissionStatusPending {
		return fmt.Errorf("can only recalculate pending commissions")
	}

	// Update fee and/or rate if provided
	if newRegistrationFee != nil {
		if *newRegistrationFee <= 0 {
			return fmt.Errorf("registration fee must be greater than 0")
		}
		commission.RegistrationFee = *newRegistrationFee
	}

	if newCommissionRate != nil {
		if *newCommissionRate < 0 || *newCommissionRate > 100 {
			return fmt.Errorf("commission rate must be between 0 and 100")
		}
		commission.CommissionRate = *newCommissionRate
	}

	// Recalculate amount
	commission.CommissionAmount = s.calculateCommissionAmount(commission.RegistrationFee, commission.CommissionRate)

	if err := s.repo.Update(schema, commission); err != nil {
		return fmt.Errorf("failed to recalculate commission: %w", err)
	}

	// Publish event
	if err := s.publishEvent("commission.recalculated", commission); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	return nil
}

// calculateCommissionAmount calculates commission amount from fee and rate
func (s *CommissionService) calculateCommissionAmount(registrationFee, commissionRate float64) float64 {
	// Calculate commission: fee * (rate / 100)
	// Round to 2 decimal places
	amount := registrationFee * (commissionRate / 100)
	return float64(int(amount*100+0.5)) / 100
}

// publishEvent publishes an event to RabbitMQ
func (s *CommissionService) publishEvent(eventType string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.rabbitCh.PublishWithContext(
		ctx,
		"comply360.commissions", // exchange
		eventType,               // routing key
		false,                   // mandatory
		false,                   // immediate
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
func (s *CommissionService) Close() error {
	if s.rabbitCh != nil {
		return s.rabbitCh.Close()
	}
	return nil
}
