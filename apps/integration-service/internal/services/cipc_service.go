package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/comply360/integration-service/internal/adapters"
	amqp "github.com/rabbitmq/amqp091-go"
)

// CIPCService handles CIPC integration operations
type CIPCService struct {
	client  *adapters.CIPCClient
	channel *amqp.Channel
}

// NewCIPCService creates a new CIPC service
func NewCIPCService(client *adapters.CIPCClient, rabbitConn *amqp.Connection) (*CIPCService, error) {
	channel, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	// Declare CIPC event queue
	_, err = channel.QueueDeclare(
		"cipc_events", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare cipc_events queue: %w", err)
	}

	log.Println("CIPC Service initialized")
	return &CIPCService{
		client:  client,
		channel: channel,
	}, nil
}

// SearchCompany searches for companies by name
func (s *CIPCService) SearchCompany(companyName string) ([]adapters.CIPCCompanySearchResult, error) {
	log.Printf("Service: Searching CIPC for company: %s", companyName)

	results, err := s.client.SearchCompanyByName(companyName)
	if err != nil {
		return nil, fmt.Errorf("failed to search CIPC: %w", err)
	}

	// Publish event
	s.publishEvent("company_searched", map[string]interface{}{
		"company_name":   companyName,
		"results_count": len(results),
	})

	return results, nil
}

// GetCompanyDetails retrieves detailed company information
func (s *CIPCService) GetCompanyDetails(registrationNumber string) (*adapters.CIPCCompanyDetails, error) {
	log.Printf("Service: Getting CIPC company details for: %s", registrationNumber)

	details, err := s.client.GetCompanyDetails(registrationNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get company details: %w", err)
	}

	// Publish event
	s.publishEvent("company_details_retrieved", map[string]interface{}{
		"registration_number": registrationNumber,
		"company_name":       details.CompanyName,
		"status":             details.Status,
	})

	return details, nil
}

// ValidateCompany validates a company's registration
func (s *CIPCService) ValidateCompany(registrationNumber string) (*adapters.CIPCValidationResult, error) {
	log.Printf("Service: Validating CIPC registration: %s", registrationNumber)

	result, err := s.client.ValidateCompany(registrationNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to validate company: %w", err)
	}

	// Publish event
	s.publishEvent("company_validated", map[string]interface{}{
		"registration_number": registrationNumber,
		"valid":              result.Valid,
		"status":             result.Status,
	})

	return result, nil
}

// CheckStatus checks a company's current status
func (s *CIPCService) CheckStatus(registrationNumber string) (string, error) {
	log.Printf("Service: Checking CIPC status for: %s", registrationNumber)

	status, err := s.client.CheckCompanyStatus(registrationNumber)
	if err != nil {
		return "", fmt.Errorf("failed to check company status: %w", err)
	}

	// Publish event
	s.publishEvent("company_status_checked", map[string]interface{}{
		"registration_number": registrationNumber,
		"status":             status,
	})

	return status, nil
}

// publishEvent publishes an event to RabbitMQ
func (s *CIPCService) publishEvent(eventType string, data map[string]interface{}) {
	event := map[string]interface{}{
		"event_type": eventType,
		"timestamp":  fmt.Sprintf("%d", amqp.Publishing{}.Timestamp.Unix()),
		"data":       data,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	err = s.channel.Publish(
		"",            // exchange
		"cipc_events", // routing key (queue name)
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish event: %v", err)
	} else {
		log.Printf("Published CIPC event: %s", eventType)
	}
}

// Close closes the RabbitMQ channel
func (s *CIPCService) Close() error {
	if s.channel != nil {
		return s.channel.Close()
	}
	return nil
}
