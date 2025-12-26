package services

import (
	"fmt"
	"log"
)

// SMSService handles sending SMS messages
type SMSService struct {
	apiKey     string
	apiSecret  string
	fromNumber string
	provider   string // twilio, africastalking, etc.
}

// SMSMessage represents an SMS to be sent
type SMSMessage struct {
	To      string
	Message string
}

// NewSMSService creates a new SMS service
func NewSMSService(apiKey, apiSecret, fromNumber, provider string) *SMSService {
	return &SMSService{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		fromNumber: fromNumber,
		provider:   provider,
	}
}

// SendSMS sends an SMS message
func (s *SMSService) SendSMS(msg SMSMessage) error {
	// TODO: Implement actual SMS sending via provider API
	// For now, just log the SMS (development mode)
	log.Printf("[SMS] To: %s | Message: %s\n", msg.To, msg.Message)

	// In production, you would integrate with an SMS provider
	// Example for Twilio:
	/*
	accountSid := s.apiKey
	authToken := s.apiSecret

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	msgData := url.Values{}
	msgData.Set("To", msg.To)
	msgData.Set("From", s.fromNumber)
	msgData.Set("Body", msg.Message)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(msgData.Encode()))
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("SMS provider returned error status: %d", resp.StatusCode)
	}
	*/

	return nil
}

// SendRegistrationStatusSMS sends SMS notification for registration status changes
func (s *SMSService) SendRegistrationStatusSMS(phoneNumber, clientName, companyName, status string) error {
	var message string

	switch status {
	case "submitted":
		message = fmt.Sprintf("Hi %s, your registration for %s has been submitted. Ref: We'll review it soon!", clientName, companyName)
	case "approved":
		message = fmt.Sprintf("Hi %s, great news! Your registration for %s has been approved!", clientName, companyName)
	case "rejected":
		message = fmt.Sprintf("Hi %s, your registration for %s requires changes. Please check your email for details.", clientName, companyName)
	default:
		message = fmt.Sprintf("Hi %s, registration status for %s: %s", clientName, companyName, status)
	}

	msg := SMSMessage{
		To:      phoneNumber,
		Message: message,
	}

	return s.SendSMS(msg)
}

// SendDocumentVerifiedSMS sends SMS notification when document is verified
func (s *SMSService) SendDocumentVerifiedSMS(phoneNumber, clientName, documentType string) error {
	msg := SMSMessage{
		To:      phoneNumber,
		Message: fmt.Sprintf("Hi %s, your %s has been verified successfully!", clientName, documentType),
	}

	return s.SendSMS(msg)
}

// SendCommissionPaidSMS sends SMS notification when commission is paid
func (s *SMSService) SendCommissionPaidSMS(phoneNumber, agentName string, amount float64, currency string) error {
	msg := SMSMessage{
		To:      phoneNumber,
		Message: fmt.Sprintf("Hi %s, your commission of %s %.2f has been paid. Check your account!", agentName, currency, amount),
	}

	return s.SendSMS(msg)
}
