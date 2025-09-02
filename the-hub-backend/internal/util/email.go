package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/config"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type EmailTemplate struct {
	Subject string
	Body    string
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
		FromName:     os.Getenv("FROM_NAME"),
	}
}

// GenerateResetToken generates a secure random token for password reset
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SendPasswordResetEmail sends a password reset email to the user
func (es *EmailService) SendPasswordResetEmail(toEmail, resetToken string) error {
	if es.SMTPHost == "" || es.SMTPPort == "" {
		config.Logger.Warn("SMTP not configured, skipping email send")
		return nil // Don't fail if email is not configured
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), resetToken)

	template := EmailTemplate{
		Subject: "Reset Your Password - The Hub",
		Body: fmt.Sprintf(`Hello,

You have requested to reset your password for The Hub.

Please click the following link to reset your password:
%s

This link will expire in 1 hour for security reasons.

If you did not request this password reset, please ignore this email.

Best regards,
The Hub Team

---
This is an automated message. Please do not reply to this email.`, resetURL),
	}

	return es.sendEmail(toEmail, template)
}

// sendEmail sends an email using SMTP
func (es *EmailService) sendEmail(toEmail string, template EmailTemplate) error {
	// Set up authentication information
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)

	// Construct the email
	to := []string{toEmail}

	// Email message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", toEmail, template.Subject, template.Body))

	// Send email
	addr := fmt.Sprintf("%s:%s", es.SMTPHost, es.SMTPPort)
	err := smtp.SendMail(addr, auth, es.FromEmail, to, msg)
	if err != nil {
		config.Logger.Errorf("Failed to send email to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	config.Logger.Infof("Password reset email sent successfully to %s", toEmail)
	return nil
}

// ValidateEmailFormat performs basic email validation
func ValidateEmailFormat(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]
	if len(local) == 0 || len(domain) == 0 {
		return false
	}

	// Check for basic domain format
	domainParts := strings.Split(domain, ".")
	return len(domainParts) >= 2 && len(domainParts[len(domainParts)-1]) >= 2
}
