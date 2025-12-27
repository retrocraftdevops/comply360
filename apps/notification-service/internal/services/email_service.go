package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
}

// EmailMessage represents an email to be sent
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// NewEmailService creates a new email service
func NewEmailService(smtpHost string, smtpPort int, smtpUsername, smtpPassword, fromEmail, fromName string) *EmailService {
	return &EmailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
		fromName:     fromName,
	}
}

// SendEmail sends an email via SMTP
func (s *EmailService) SendEmail(msg EmailMessage) error {
	// Development mode: Just log if SMTP not configured
	if s.smtpHost == "" || s.smtpHost == "localhost" {
		log.Printf("[EMAIL - DEV MODE] To: %v | Subject: %s | Body: %s\n", msg.To, msg.Subject, msg.Body)
		return nil
	}

	// Production mode: Actually send email
	log.Printf("[EMAIL] Sending to: %v | Subject: %s", msg.To, msg.Subject)

	// Build email headers
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(msg.To, ", ")
	headers["Subject"] = msg.Subject
	headers["MIME-Version"] = "1.0"

	if msg.IsHTML {
		headers["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
	}

	// Build message
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(msg.Body)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	// Try TLS connection first (port 465)
	if s.smtpPort == 465 {
		return s.sendViaTLS(addr, msg.To, []byte(message.String()))
	}

	// Use STARTTLS (port 587 or 25)
	return s.sendViaSTARTTLS(addr, msg.To, []byte(message.String()))
}

// sendViaTLS sends email using TLS connection (port 465)
func (s *EmailService) sendViaTLS(addr string, to []string, message []byte) error {
	// TLS config
	tlsConfig := &tls.Config{
		ServerName: s.smtpHost,
	}

	// Connect with TLS
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect via TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if s.smtpUsername != "" {
		auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	// Send email
	if err := client.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT command failed for %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Printf("[EMAIL] Successfully sent to: %v", to)
	return nil
}

// sendViaSTARTTLS sends email using STARTTLS (port 587 or 25)
func (s *EmailService) sendViaSTARTTLS(addr string, to []string, message []byte) error {
	// Connect to SMTP server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// STARTTLS if supported
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: s.smtpHost,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			log.Printf("[EMAIL] Warning: STARTTLS failed: %v", err)
		}
	}

	// Authenticate
	if s.smtpUsername != "" {
		auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	// Send email
	if err := client.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT command failed for %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	log.Printf("[EMAIL] Successfully sent to: %v", to)
	return nil
}

// SendRegistrationCreatedEmail sends notification when registration is created
func (s *EmailService) SendRegistrationCreatedEmail(toEmail, clientName, companyName string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Registration Created - " + companyName,
		Body: fmt.Sprintf(`Dear %s,

Your company registration for %s has been successfully created and is now in draft status.

Next steps:
1. Complete all required information
2. Upload supporting documents
3. Submit for review

You can track the progress of your registration in your dashboard.

Best regards,
Comply360 Team`, clientName, companyName),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendRegistrationSubmittedEmail sends notification when registration is submitted
func (s *EmailService) SendRegistrationSubmittedEmail(toEmail, clientName, companyName, registrationNumber string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Registration Submitted - " + companyName,
		Body: fmt.Sprintf(`Dear %s,

Your company registration for %s has been successfully submitted for review.

Registration Reference: %s

Our team will review your application and supporting documents. You will receive an email notification once the review is complete.

Typical review time: 2-5 business days

Best regards,
Comply360 Team`, clientName, companyName, registrationNumber),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendRegistrationApprovedEmail sends notification when registration is approved
func (s *EmailService) SendRegistrationApprovedEmail(toEmail, clientName, companyName, registrationNumber string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Registration Approved - " + companyName,
		Body: fmt.Sprintf(`Dear %s,

Congratulations! Your company registration for %s has been approved.

Registration Number: %s

The registration documents will be processed and you will receive the official registration certificate shortly.

Best regards,
Comply360 Team`, clientName, companyName, registrationNumber),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendRegistrationRejectedEmail sends notification when registration is rejected
func (s *EmailService) SendRegistrationRejectedEmail(toEmail, clientName, companyName, reason string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Registration Rejected - " + companyName,
		Body: fmt.Sprintf(`Dear %s,

Unfortunately, your company registration for %s has been rejected.

Reason: %s

You can review the feedback and resubmit your application after addressing the issues mentioned above.

If you have any questions, please contact our support team.

Best regards,
Comply360 Team`, clientName, companyName, reason),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendDocumentUploadedEmail sends notification when document is uploaded
func (s *EmailService) SendDocumentUploadedEmail(toEmail, clientName, documentType, fileName string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Document Uploaded - " + documentType,
		Body: fmt.Sprintf(`Dear %s,

Your document has been successfully uploaded.

Document Type: %s
File Name: %s

The document will be reviewed by our team shortly.

Best regards,
Comply360 Team`, clientName, documentType, fileName),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendDocumentVerifiedEmail sends notification when document is verified
func (s *EmailService) SendDocumentVerifiedEmail(toEmail, clientName, documentType string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Document Verified - " + documentType,
		Body: fmt.Sprintf(`Dear %s,

Your document has been successfully verified.

Document Type: %s

Status: Verified ✓

Best regards,
Comply360 Team`, clientName, documentType),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendCommissionApprovedEmail sends notification when commission is approved
func (s *EmailService) SendCommissionApprovedEmail(toEmail, agentName string, amount float64, currency string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Commission Approved",
		Body: fmt.Sprintf(`Dear %s,

Your commission has been approved.

Amount: %s %.2f

The payment will be processed according to the payment schedule.

Best regards,
Comply360 Team`, agentName, currency, amount),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}

// SendCommissionPaidEmail sends notification when commission is paid
func (s *EmailService) SendCommissionPaidEmail(toEmail, agentName string, amount float64, currency, paymentReference string) error {
	msg := EmailMessage{
		To:      []string{toEmail},
		Subject: "Commission Paid",
		Body: fmt.Sprintf(`Dear %s,

Your commission has been paid.

Amount: %s %.2f
Payment Reference: %s

Please check your account for the payment.

Best regards,
Comply360 Team`, agentName, currency, amount, paymentReference),
		IsHTML: false,
	}

	return s.SendEmail(msg)
}
