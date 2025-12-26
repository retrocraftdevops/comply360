package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/comply360/integration-service/internal/services"
	"github.com/comply360/shared/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

// OdooSyncConsumer handles syncing events to Odoo ERP
type OdooSyncConsumer struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	odooService *services.OdooService
}

// NewOdooSyncConsumer creates a new Odoo sync consumer
func NewOdooSyncConsumer(rabbitConn *amqp.Connection, odooService *services.OdooService) (*OdooSyncConsumer, error) {
	ch, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &OdooSyncConsumer{
		conn:        rabbitConn,
		channel:     ch,
		odooService: odooService,
	}, nil
}

// Start starts consuming events for Odoo sync
func (c *OdooSyncConsumer) Start() error {
	// Setup queues and bindings
	if err := c.setupQueues(); err != nil {
		return fmt.Errorf("failed to setup queues: %w", err)
	}

	// Start consuming registration events
	go c.consumeRegistrationEvents()

	// Start consuming commission events
	go c.consumeCommissionEvents()

	log.Println("Odoo sync consumer started successfully")
	return nil
}

// setupQueues sets up RabbitMQ queues and bindings for Odoo sync
func (c *OdooSyncConsumer) setupQueues() error {
	// Declare queues for Odoo sync
	queues := []string{
		"comply360.odoo.registration",
		"comply360.odoo.commission",
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

	// Bind registration events for Odoo sync
	registrationBindings := []string{
		"registration.status.submitted", // Create lead when submitted
		"registration.status.approved",  // Update lead to won
	}

	for _, routingKey := range registrationBindings {
		err := c.channel.QueueBind(
			"comply360.odoo.registration", // queue name
			routingKey,                     // routing key
			"comply360.registrations",      // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue for %s: %w", routingKey, err)
		}
	}

	// Bind commission events for Odoo sync
	commissionBindings := []string{
		"commission.approved", // Create invoice when approved
		"commission.paid",     // Mark invoice as paid
	}

	for _, routingKey := range commissionBindings {
		err := c.channel.QueueBind(
			"comply360.odoo.commission", // queue name
			routingKey,                   // routing key
			"comply360.commissions",      // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue for %s: %w", routingKey, err)
		}
	}

	log.Println("Odoo sync queues and bindings configured successfully")
	return nil
}

// consumeRegistrationEvents consumes registration events for Odoo sync
func (c *OdooSyncConsumer) consumeRegistrationEvents() {
	msgs, err := c.channel.Consume(
		"comply360.odoo.registration", // queue
		"odoo-sync-registration",       // consumer
		false,                          // auto-ack
		false,                          // exclusive
		false,                          // no-local
		false,                          // no-wait
		nil,                            // args
	)
	if err != nil {
		log.Printf("Failed to consume registration events for Odoo sync: %v", err)
		return
	}

	for msg := range msgs {
		c.handleRegistrationEvent(msg)
	}
}

// consumeCommissionEvents consumes commission events for Odoo sync
func (c *OdooSyncConsumer) consumeCommissionEvents() {
	msgs, err := c.channel.Consume(
		"comply360.odoo.commission", // queue
		"odoo-sync-commission",       // consumer
		false,                        // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		log.Printf("Failed to consume commission events for Odoo sync: %v", err)
		return
	}

	for msg := range msgs {
		c.handleCommissionEvent(msg)
	}
}

// handleRegistrationEvent syncs registration events to Odoo
func (c *OdooSyncConsumer) handleRegistrationEvent(msg amqp.Delivery) {
	routingKey := msg.RoutingKey
	log.Printf("[ODOO SYNC] Received registration event: %s", routingKey)

	var registration models.Registration
	if err := json.Unmarshal(msg.Body, &registration); err != nil {
		log.Printf("[ODOO SYNC] Failed to unmarshal registration event: %v", err)
		msg.Nack(false, false)
		return
	}

	switch routingKey {
	case "registration.status.submitted":
		// Sync registration to Odoo when submitted
		log.Printf("[ODOO SYNC] Syncing registration to Odoo: %s", registration.CompanyName)

		// In production, you would:
		// 1. Convert models.Registration to integration service Registration model
		// 2. Call odooService.CreateLeadFromRegistration()
		// 3. Update registration with returned Odoo lead ID via API call back to registration service

		log.Printf("[ODOO SYNC] Would create Odoo lead for: %s (Type: %s, Jurisdiction: %s)",
			registration.CompanyName, registration.RegistrationType, registration.Jurisdiction)

		// Mock the lead creation
		log.Printf("[ODOO SYNC] ✓ [DEMO] Would create Odoo lead for registration")

	case "registration.status.approved":
		// Sync approval to Odoo when registration is approved
		log.Printf("[ODOO SYNC] Syncing registration approval to Odoo: %s", registration.CompanyName)

		// In production, you would:
		// 1. Find the Odoo lead by registration ID
		// 2. Call odooService.SyncRegistrationStatus() to update status
		// 3. Convert lead to customer if needed

		log.Printf("[ODOO SYNC] Would update Odoo lead status to 'won' for: %s", registration.CompanyName)
		log.Printf("[ODOO SYNC] ✓ [DEMO] Would convert Odoo lead to customer")
	}

	msg.Ack(false)
	log.Printf("[ODOO SYNC] Successfully processed registration event: %s", routingKey)
}

// handleCommissionEvent syncs commission events to Odoo
func (c *OdooSyncConsumer) handleCommissionEvent(msg amqp.Delivery) {
	routingKey := msg.RoutingKey
	log.Printf("[ODOO SYNC] Received commission event: %s", routingKey)

	var commission models.Commission
	if err := json.Unmarshal(msg.Body, &commission); err != nil {
		log.Printf("[ODOO SYNC] Failed to unmarshal commission event: %v", err)
		msg.Nack(false, false)
		return
	}

	switch routingKey {
	case "commission.approved":
		// Sync commission approval to Odoo
		log.Printf("[ODOO SYNC] Syncing commission approval to Odoo: %s %.2f",
			commission.Currency, commission.CommissionAmount)

		// In production, you would:
		// 1. Convert models.Commission to integration service Commission model
		// 2. Call odooService.CreateCommission() to create commission record
		// 3. Or create an invoice for the commission
		// 4. Update commission with returned Odoo record ID via API call

		log.Printf("[ODOO SYNC] Would create Odoo commission/invoice for agent: %s", commission.AgentID)
		log.Printf("[ODOO SYNC] ✓ [DEMO] Would create Odoo commission record")

	case "commission.paid":
		// Sync commission payment to Odoo
		log.Printf("[ODOO SYNC] Syncing commission payment to Odoo")

		paymentRef := "N/A"
		if commission.PaymentReference != nil {
			paymentRef = *commission.PaymentReference
		}

		// In production, you would:
		// 1. Find the related Odoo invoice/commission record
		// 2. Register payment against the invoice
		// 3. Mark as paid in Odoo

		log.Printf("[ODOO SYNC] Would register payment in Odoo: %s %.2f (Ref: %s)",
			commission.Currency, commission.CommissionAmount, paymentRef)
		log.Printf("[ODOO SYNC] ✓ [DEMO] Would mark Odoo commission as paid")
	}

	msg.Ack(false)
	log.Printf("[ODOO SYNC] Successfully processed commission event: %s", routingKey)
}

// Close closes the consumer
func (c *OdooSyncConsumer) Close() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}
