package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(email, token string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)
	message := []byte("Subject: Email Verification\n\nClick the link to verify your email: " + verificationLink)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUser, []string{email}, message)
	if err != nil {
		return err
	}
	return nil
}
