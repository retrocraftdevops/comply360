package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/comply360/notification-service/internal/services"
	"github.com/comply360/shared/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EventConsumer handles consuming events from RabbitMQ
type EventConsumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	emailService *services.EmailService
	smsService   *services.SMSService
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(rabbitConn *amqp.Connection, emailService *services.EmailService, smsService *services.SMSService) (*EventConsumer, error) {
	ch, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &EventConsumer{
		conn:         rabbitConn,
		channel:      ch,
		emailService: emailService,
		smsService:   smsService,
	}, nil
}

// Start starts consuming events from all exchanges
func (c *EventConsumer) Start() error {
	// Setup queues and bindings
	if err := c.setupQueues(); err != nil {
		return fmt.Errorf("failed to setup queues: %w", err)
	}

	// Start consuming registration events
	go c.consumeRegistrationEvents()

	// Start consuming document events
	go c.consumeDocumentEvents()

	// Start consuming commission events
	go c.consumeCommissionEvents()

	log.Println("Event consumer started successfully")
	return nil
}

// setupQueues sets up RabbitMQ queues and bindings
func (c *EventConsumer) setupQueues() error {
	// Declare queues
	queues := []string{
		"comply360.notifications.registration",
		"comply360.notifications.document",
		"comply360.notifications.commission",
	}

	for _, queueName := range queues {
		_, err := c.channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	// Bind registration events
	registrationBindings := []string{
		"registration.created",
		"registration.status.submitted",
		"registration.status.approved",
		"registration.status.rejected",
	}

	for _, routingKey := range registrationBindings {
		err := c.channel.QueueBind(
			"comply360.notifications.registration", // queue name
			routingKey,                              // routing key
			"comply360.registrations",               // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue for %s: %w", routingKey, err)
		}
	}

	// Bind document events
	documentBindings := []string{
		"document.uploaded",
		"document.verified",
	}

	for _, routingKey := range documentBindings {
		err := c.channel.QueueBind(
			"comply360.notifications.document", // queue name
			routingKey,                          // routing key
			"comply360.documents",               // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue for %s: %w", routingKey, err)
		}
	}

	// Bind commission events
	commissionBindings := []string{
		"commission.approved",
		"commission.paid",
	}

	for _, routingKey := range commissionBindings {
		err := c.channel.QueueBind(
			"comply360.notifications.commission", // queue name
			routingKey,                            // routing key
			"comply360.commissions",               // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue for %s: %w", routingKey, err)
		}
	}

	log.Println("RabbitMQ queues and bindings configured successfully")
	return nil
}

// consumeRegistrationEvents consumes registration events
func (c *EventConsumer) consumeRegistrationEvents() {
	msgs, err := c.channel.Consume(
		"comply360.notifications.registration", // queue
		"",                                      // consumer
		false,                                   // auto-ack
		false,                                   // exclusive
		false,                                   // no-local
		false,                                   // no-wait
		nil,                                     // args
	)
	if err != nil {
		log.Printf("Failed to consume registration events: %v", err)
		return
	}

	for msg := range msgs {
		c.handleRegistrationEvent(msg)
	}
}

// consumeDocumentEvents consumes document events
func (c *EventConsumer) consumeDocumentEvents() {
	msgs, err := c.channel.Consume(
		"comply360.notifications.document", // queue
		"",                                  // consumer
		false,                               // auto-ack
		false,                               // exclusive
		false,                               // no-local
		false,                               // no-wait
		nil,                                 // args
	)
	if err != nil {
		log.Printf("Failed to consume document events: %v", err)
		return
	}

	for msg := range msgs {
		c.handleDocumentEvent(msg)
	}
}

// consumeCommissionEvents consumes commission events
func (c *EventConsumer) consumeCommissionEvents() {
	msgs, err := c.channel.Consume(
		"comply360.notifications.commission", // queue
		"",                                    // consumer
		false,                                 // auto-ack
		false,                                 // exclusive
		false,                                 // no-local
		false,                                 // no-wait
		nil,                                   // args
	)
	if err != nil {
		log.Printf("Failed to consume commission events: %v", err)
		return
	}

	for msg := range msgs {
		c.handleCommissionEvent(msg)
	}
}

// handleRegistrationEvent handles registration events
func (c *EventConsumer) handleRegistrationEvent(msg amqp.Delivery) {
	routingKey := msg.RoutingKey
	log.Printf("Received registration event: %s", routingKey)

	var registration models.Registration
	if err := json.Unmarshal(msg.Body, &registration); err != nil {
		log.Printf("Failed to unmarshal registration event: %v", err)
		msg.Nack(false, false) // Don't requeue malformed messages
		return
	}

	// TODO: Fetch client details from database to get email/phone
	// For now, using placeholder values
	clientEmail := "client@example.com"
	clientName := "Client Name"

	var err error
	switch routingKey {
	case "registration.created":
		err = c.emailService.SendRegistrationCreatedEmail(clientEmail, clientName, registration.CompanyName)
	case "registration.status.submitted":
		err = c.emailService.SendRegistrationSubmittedEmail(clientEmail, clientName, registration.CompanyName, registration.ID.String())
		// Also send SMS
		c.smsService.SendRegistrationStatusSMS("+1234567890", clientName, registration.CompanyName, "submitted")
	case "registration.status.approved":
		regNumber := registration.ID.String()
		if registration.RegistrationNumber != nil {
			regNumber = *registration.RegistrationNumber
		}
		err = c.emailService.SendRegistrationApprovedEmail(clientEmail, clientName, registration.CompanyName, regNumber)
		c.smsService.SendRegistrationStatusSMS("+1234567890", clientName, registration.CompanyName, "approved")
	case "registration.status.rejected":
		reason := "Please review and resubmit"
		if registration.RejectionReason != nil {
			reason = *registration.RejectionReason
		}
		err = c.emailService.SendRegistrationRejectedEmail(clientEmail, clientName, registration.CompanyName, reason)
		c.smsService.SendRegistrationStatusSMS("+1234567890", clientName, registration.CompanyName, "rejected")
	}

	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		msg.Nack(false, true) // Requeue on error
		return
	}

	msg.Ack(false)
	log.Printf("Successfully processed registration event: %s", routingKey)
}

// handleDocumentEvent handles document events
func (c *EventConsumer) handleDocumentEvent(msg amqp.Delivery) {
	routingKey := msg.RoutingKey
	log.Printf("Received document event: %s", routingKey)

	var document models.Document
	if err := json.Unmarshal(msg.Body, &document); err != nil {
		log.Printf("Failed to unmarshal document event: %v", err)
		msg.Nack(false, false)
		return
	}

	// TODO: Fetch client details from database
	clientEmail := "client@example.com"
	clientName := "Client Name"

	var err error
	switch routingKey {
	case "document.uploaded":
		err = c.emailService.SendDocumentUploadedEmail(clientEmail, clientName, document.DocumentType, document.FileName)
	case "document.verified":
		err = c.emailService.SendDocumentVerifiedEmail(clientEmail, clientName, document.DocumentType)
		c.smsService.SendDocumentVerifiedSMS("+1234567890", clientName, document.DocumentType)
	}

	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	log.Printf("Successfully processed document event: %s", routingKey)
}

// handleCommissionEvent handles commission events
func (c *EventConsumer) handleCommissionEvent(msg amqp.Delivery) {
	routingKey := msg.RoutingKey
	log.Printf("Received commission event: %s", routingKey)

	var commission models.Commission
	if err := json.Unmarshal(msg.Body, &commission); err != nil {
		log.Printf("Failed to unmarshal commission event: %v", err)
		msg.Nack(false, false)
		return
	}

	// TODO: Fetch agent details from database
	agentEmail := "agent@example.com"
	agentName := "Agent Name"

	var err error
	switch routingKey {
	case "commission.approved":
		err = c.emailService.SendCommissionApprovedEmail(agentEmail, agentName, commission.CommissionAmount, commission.Currency)
	case "commission.paid":
		paymentRef := "N/A"
		if commission.PaymentReference != nil {
			paymentRef = *commission.PaymentReference
		}
		err = c.emailService.SendCommissionPaidEmail(agentEmail, agentName, commission.CommissionAmount, commission.Currency, paymentRef)
		c.smsService.SendCommissionPaidSMS("+1234567890", agentName, commission.CommissionAmount, commission.Currency)
	}

	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	log.Printf("Successfully processed commission event: %s", routingKey)
}

// Close closes the consumer
func (c *EventConsumer) Close() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}
