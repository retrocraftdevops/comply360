package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// CIPCConfig holds CIPC API connection configuration
type CIPCConfig struct {
	BaseURL  string
	APIKey   string
	Username string
	Password string
}

// CIPCClient handles REST API communication with CIPC
type CIPCClient struct {
	config     *CIPCConfig
	httpClient *http.Client
	token      string
	tokenExpiry time.Time
}

// CIPCCompanySearchResult represents a company search result from CIPC
type CIPCCompanySearchResult struct {
	RegistrationNumber string `json:"registration_number"`
	CompanyName        string `json:"company_name"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	RegistrationDate   string `json:"registration_date"`
}

// CIPCCompanyDetails represents detailed company information from CIPC
type CIPCCompanyDetails struct {
	RegistrationNumber string   `json:"registration_number"`
	CompanyName        string   `json:"company_name"`
	Status             string   `json:"status"`
	Type               string   `json:"type"`
	RegistrationDate   string   `json:"registration_date"`
	BusinessAddress    string   `json:"business_address"`
	PostalAddress      string   `json:"postal_address"`
	Directors          []string `json:"directors"`
	EntityNumber       string   `json:"entity_number"`
}

// CIPCValidationResult represents the result of a company validation
type CIPCValidationResult struct {
	Valid              bool   `json:"valid"`
	RegistrationNumber string `json:"registration_number"`
	CompanyName        string `json:"company_name"`
	Status             string `json:"status"`
	Message            string `json:"message"`
}

// NewCIPCClient creates a new CIPC API client
func NewCIPCClient(config *CIPCConfig) (*CIPCClient, error) {
	client := &CIPCClient{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Note: In production, this would authenticate with CIPC API
	// For development, we'll simulate the authentication
	log.Printf("CIPC Client initialized for: %s", config.BaseURL)
	return client, nil
}

// authenticate gets an authentication token from CIPC
func (c *CIPCClient) authenticate() error {
	// In production, this would call CIPC's authentication endpoint
	// For now, we'll simulate it
	c.token = "simulated_cipc_token"
	c.tokenExpiry = time.Now().Add(1 * time.Hour)
	return nil
}

// ensureAuthenticated checks if token is valid and refreshes if needed
func (c *CIPCClient) ensureAuthenticated() error {
	if c.token == "" || time.Now().After(c.tokenExpiry) {
		return c.authenticate()
	}
	return nil
}

// SearchCompanyByName searches for companies by name
func (c *CIPCClient) SearchCompanyByName(companyName string) ([]CIPCCompanySearchResult, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	// In production, this would call CIPC API
	// For development, we'll return simulated results
	log.Printf("CIPC: Searching for company: %s", companyName)

	// Simulate API response
	results := []CIPCCompanySearchResult{
		{
			RegistrationNumber: "2024/999999/07",
			CompanyName:        companyName,
			Status:             "Active",
			Type:               "Private Company",
			RegistrationDate:   "2024-01-15",
		},
	}

	return results, nil
}

// GetCompanyDetails retrieves detailed information about a company
func (c *CIPCClient) GetCompanyDetails(registrationNumber string) (*CIPCCompanyDetails, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	log.Printf("CIPC: Getting details for registration: %s", registrationNumber)

	// In production, this would call CIPC API
	// For development, we'll return simulated details
	details := &CIPCCompanyDetails{
		RegistrationNumber: registrationNumber,
		CompanyName:        "Sample Company (Pty) Ltd",
		Status:             "Active",
		Type:               "Private Company",
		RegistrationDate:   "2024-01-15",
		BusinessAddress:    "123 Main Street, Johannesburg, 2000",
		PostalAddress:      "PO Box 12345, Johannesburg, 2000",
		Directors:          []string{"John Doe", "Jane Smith"},
		EntityNumber:       "K2024999999",
	}

	return details, nil
}

// ValidateCompany validates a company's registration status
func (c *CIPCClient) ValidateCompany(registrationNumber string) (*CIPCValidationResult, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	log.Printf("CIPC: Validating registration: %s", registrationNumber)

	// In production, this would call CIPC validation API
	// For development, we'll return simulated validation
	result := &CIPCValidationResult{
		Valid:              true,
		RegistrationNumber: registrationNumber,
		CompanyName:        "Sample Company (Pty) Ltd",
		Status:             "Active",
		Message:            "Company registration is valid and active",
	}

	return result, nil
}

// CheckCompanyStatus checks the current status of a company
func (c *CIPCClient) CheckCompanyStatus(registrationNumber string) (string, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return "", err
	}

	log.Printf("CIPC: Checking status for registration: %s", registrationNumber)

	// In production, this would call CIPC status API
	// Possible statuses: Active, Deregistered, In Liquidation, etc.
	return "Active", nil
}

// makeRequest makes an HTTP request to CIPC API (helper method for production use)
func (c *CIPCClient) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.config.BaseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.config.APIKey != "" {
		req.Header.Set("X-API-Key", c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
