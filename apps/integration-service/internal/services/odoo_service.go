package services

import (
	"fmt"
	"log"
	"time"

	"github.com/comply360/integration-service/internal/adapters"
	"github.com/comply360/integration-service/internal/models"
	"github.com/google/uuid"
)

// OdooService handles Odoo 19 ERP integration logic
type OdooService struct {
	client *adapters.OdooClient
}

// NewOdooService creates a new Odoo service
func NewOdooService(client *adapters.OdooClient) *OdooService {
	return &OdooService{
		client: client,
	}
}

// CreateLeadFromRegistration creates a CRM lead in Odoo from a Comply360 registration
func (s *OdooService) CreateLeadFromRegistration(registration *models.Registration) (int, error) {
	log.Printf("Creating Odoo lead from registration: %s", registration.ID)

	leadValues := map[string]interface{}{
		"name":         registration.CompanyName,
		"contact_name": registration.ContactPerson,
		"email_from":   registration.Email,
		"phone":        registration.Phone,
		"type":         "opportunity",
		"description":  fmt.Sprintf("Registration Type: %s\nRegistration Number: %s\nCreated from Comply360", registration.RegistrationType, registration.RegistrationNumber),
	}

	// Add address if available
	if registration.Address != "" {
		leadValues["street"] = registration.Address
	}
	if registration.City != "" {
		leadValues["city"] = registration.City
	}
	if registration.Country != "" {
		leadValues["country_id"] = s.getCountryID(registration.Country)
	}

	// Add custom field for Comply360 registration ID
	leadValues["x_comply360_registration_id"] = registration.ID.String()

	// Set probability based on status
	switch registration.Status {
	case "submitted":
		leadValues["probability"] = 10
	case "under_review":
		leadValues["probability"] = 30
	case "approved":
		leadValues["probability"] = 100
	case "rejected":
		leadValues["probability"] = 0
	default:
		leadValues["probability"] = 5
	}

	leadID, err := s.client.Create("crm.lead", leadValues)
	if err != nil {
		return 0, fmt.Errorf("failed to create lead: %w", err)
	}

	log.Printf("Created Odoo lead ID %d for registration %s", leadID, registration.ID)
	return leadID, nil
}

// UpdateLead updates an existing CRM lead in Odoo
func (s *OdooService) UpdateLead(leadID int, updates map[string]interface{}) error {
	log.Printf("Updating Odoo lead ID %d", leadID)

	success, err := s.client.Write("crm.lead", []int{leadID}, updates)
	if err != nil {
		return fmt.Errorf("failed to update lead: %w", err)
	}

	if !success {
		return fmt.Errorf("update lead returned false")
	}

	log.Printf("Successfully updated Odoo lead ID %d", leadID)
	return nil
}

// ConvertLeadToCustomer converts a CRM lead to a customer (res.partner)
func (s *OdooService) ConvertLeadToCustomer(leadID int) (int, error) {
	log.Printf("Converting Odoo lead ID %d to customer", leadID)

	// First, get the lead details
	leads, err := s.client.Read("crm.lead", []int{leadID}, []string{"name", "contact_name", "email_from", "phone", "street", "city", "country_id"})
	if err != nil {
		return 0, fmt.Errorf("failed to read lead: %w", err)
	}

	if len(leads) == 0 {
		return 0, fmt.Errorf("lead not found")
	}

	lead := leads[0]

	// Create partner (customer)
	partnerValues := map[string]interface{}{
		"name":       lead["name"],
		"email":      lead["email_from"],
		"phone":      lead["phone"],
		"is_company": true,
		"customer_rank": 1,
	}

	if contactName, ok := lead["contact_name"].(string); ok && contactName != "" {
		partnerValues["contact_name"] = contactName
	}
	if street, ok := lead["street"].(string); ok {
		partnerValues["street"] = street
	}
	if city, ok := lead["city"].(string); ok {
		partnerValues["city"] = city
	}
	if countryID, ok := lead["country_id"].([]interface{}); ok && len(countryID) > 0 {
		partnerValues["country_id"] = countryID[0]
	}

	partnerID, err := s.client.Create("res.partner", partnerValues)
	if err != nil {
		return 0, fmt.Errorf("failed to create partner: %w", err)
	}

	// Update lead to mark as won and link to partner
	leadUpdates := map[string]interface{}{
		"partner_id": partnerID,
		"probability": 100,
	}

	if err := s.UpdateLead(leadID, leadUpdates); err != nil {
		log.Printf("Warning: Created partner %d but failed to update lead: %v", partnerID, err)
	}

	// Mark lead as won
	_, err = s.client.Execute("crm.lead", "action_set_won", []interface{}{[]int{leadID}})
	if err != nil {
		log.Printf("Warning: Failed to mark lead as won: %v", err)
	}

	log.Printf("Converted lead %d to customer (partner) %d", leadID, partnerID)
	return partnerID, nil
}

// GetLeadByRegistrationID finds a lead by Comply360 registration ID
func (s *OdooService) GetLeadByRegistrationID(registrationID uuid.UUID) (int, error) {
	domain := []interface{}{
		[]interface{}{"x_comply360_registration_id", "=", registrationID.String()},
	}

	leadIDs, err := s.client.Search("crm.lead", domain, map[string]interface{}{"limit": 1})
	if err != nil {
		return 0, fmt.Errorf("failed to search for lead: %w", err)
	}

	if len(leadIDs) == 0 {
		return 0, fmt.Errorf("lead not found for registration %s", registrationID)
	}

	return leadIDs[0], nil
}

// CreateInvoice creates an invoice in Odoo
func (s *OdooService) CreateInvoice(invoice *models.Invoice) (int, error) {
	log.Printf("Creating Odoo invoice for customer partner ID %d", invoice.PartnerID)

	invoiceValues := map[string]interface{}{
		"partner_id":   invoice.PartnerID,
		"move_type":    "out_invoice", // Customer invoice
		"invoice_date": time.Now().Format("2006-01-02"),
		"payment_reference": invoice.Reference,
	}

	// Add invoice lines
	if len(invoice.Lines) > 0 {
		invoiceLines := make([]interface{}, len(invoice.Lines))
		for i, line := range invoice.Lines {
			invoiceLines[i] = []interface{}{
				0, 0, map[string]interface{}{
					"name":       line.Description,
					"quantity":   line.Quantity,
					"price_unit": line.UnitPrice,
				},
			}
		}
		invoiceValues["invoice_line_ids"] = invoiceLines
	}

	invoiceID, err := s.client.Create("account.move", invoiceValues)
	if err != nil {
		return 0, fmt.Errorf("failed to create invoice: %w", err)
	}

	log.Printf("Created Odoo invoice ID %d", invoiceID)
	return invoiceID, nil
}

// CreateCommission creates a commission record in Odoo
func (s *OdooService) CreateCommission(commission *models.Commission) (int, error) {
	log.Printf("Creating Odoo commission for agent %s", commission.AgentID)

	// First, find or create the agent as a partner
	agentPartnerID, err := s.getOrCreateAgent(commission.AgentID, commission.AgentName)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create agent: %w", err)
	}

	commissionValues := map[string]interface{}{
		"name":        fmt.Sprintf("Commission - %s - %s", commission.AgentName, commission.RegistrationNumber),
		"partner_id":  agentPartnerID,
		"amount":      commission.Amount,
		"date":        time.Now().Format("2006-01-02"),
		"state":       s.mapCommissionStatus(commission.Status),
		"description": fmt.Sprintf("Registration: %s\nType: %s\nComply360 Commission ID: %s", commission.RegistrationNumber, commission.RegistrationType, commission.ID),
		"x_comply360_commission_id": commission.ID,
	}

	// Note: This assumes you have a custom commission model in Odoo
	// If not, you might use hr.expense or account.move instead
	commissionID, err := s.client.Create("x_commission", commissionValues)
	if err != nil {
		return 0, fmt.Errorf("failed to create commission: %w", err)
	}

	log.Printf("Created Odoo commission ID %d", commissionID)
	return commissionID, nil
}

// SyncRegistrationStatus syncs registration status changes to Odoo lead
func (s *OdooService) SyncRegistrationStatus(registrationID uuid.UUID, status string) error {
	leadID, err := s.GetLeadByRegistrationID(registrationID)
	if err != nil {
		return fmt.Errorf("failed to find lead: %w", err)
	}

	updates := map[string]interface{}{}

	switch status {
	case "approved":
		updates["probability"] = 100
		// Convert to customer
		_, err := s.ConvertLeadToCustomer(leadID)
		if err != nil {
			log.Printf("Failed to convert lead to customer: %v", err)
		}
	case "rejected":
		updates["probability"] = 0
		updates["active"] = false
	case "under_review":
		updates["probability"] = 50
	case "submitted":
		updates["probability"] = 10
	}

	if len(updates) > 0 {
		return s.UpdateLead(leadID, updates)
	}

	return nil
}

// TestConnection tests the connection to Odoo
func (s *OdooService) TestConnection() error {
	version, err := s.client.GetServerVersion()
	if err != nil {
		return fmt.Errorf("failed to get Odoo version: %w", err)
	}

	log.Printf("Connected to Odoo version: %v", version)

	// Test access rights
	hasAccess, err := s.client.CheckAccessRights("crm.lead", "create")
	if err != nil {
		return fmt.Errorf("failed to check access rights: %w", err)
	}

	if !hasAccess {
		return fmt.Errorf("user does not have create access to crm.lead")
	}

	log.Println("Odoo connection test successful")
	return nil
}

// Helper functions

func (s *OdooService) getCountryID(countryCode string) interface{} {
	// Search for country by code
	domain := []interface{}{
		[]interface{}{"code", "=", countryCode},
	}

	countryIDs, err := s.client.Search("res.country", domain, map[string]interface{}{"limit": 1})
	if err != nil || len(countryIDs) == 0 {
		return false
	}

	return countryIDs[0]
}

func (s *OdooService) getOrCreateAgent(agentID uuid.UUID, agentName string) (int, error) {
	// Search for existing agent by Comply360 ID
	domain := []interface{}{
		[]interface{}{"x_comply360_agent_id", "=", agentID.String()},
	}

	partnerIDs, err := s.client.Search("res.partner", domain, map[string]interface{}{"limit": 1})
	if err == nil && len(partnerIDs) > 0 {
		return partnerIDs[0], nil
	}

	// Create new agent partner
	partnerValues := map[string]interface{}{
		"name":                   agentName,
		"is_company":             false,
		"supplier_rank":          1, // Mark as vendor for commission payments
		"x_comply360_agent_id":   agentID.String(),
	}

	return s.client.Create("res.partner", partnerValues)
}

func (s *OdooService) mapCommissionStatus(status string) string {
	switch status {
	case "pending":
		return "draft"
	case "approved":
		return "approved"
	case "paid":
		return "paid"
	case "rejected":
		return "refused"
	default:
		return "draft"
	}
}
