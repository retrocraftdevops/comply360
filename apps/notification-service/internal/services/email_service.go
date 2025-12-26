package services

import (
	"fmt"
	"log"
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

// SendEmail sends an email
func (s *EmailService) SendEmail(msg EmailMessage) error {
	// TODO: Implement actual email sending via SMTP
	// For now, just log the email (development mode)
	log.Printf("[EMAIL] To: %v | Subject: %s | Body: %s\n", msg.To, msg.Subject, msg.Body)

	// In production, you would use a library like gomail or net/smtp
	// Example implementation would be:
	/*
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail))
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	*/

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
