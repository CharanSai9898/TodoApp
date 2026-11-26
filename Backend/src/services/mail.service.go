package services

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"os"
	"strings"
)

// GenerateOTP produces a cryptographically secure 6-digit numeric string
func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// SendOTPEmail sends an OTP via SMTP and logs it for local development
func SendOTPEmail(toEmail, otp string) error {
	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpPort := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	smtpEmail := strings.TrimSpace(os.Getenv("SMTP_EMAIL"))
	smtpPassword := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))

	// Always print OTP in terminal for immediate testing in development
	log.Printf("==================================================")
	log.Printf(">>> [2FA OTP CODE] For %s: %s <<<", toEmail, otp)
	log.Printf("==================================================")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		log.Println("[2FA OTP INFO] SMTP credentials not fully configured in .env; use the code printed above.")
		return nil
	}

	auth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)
	subject := "Subject: Your Login Verification Code\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 500px; margin: auto; padding: 24px; border: 1px solid #e5e7eb; border-radius: 8px;">
			<h2 style="color: #2563eb; text-align: center; margin-bottom: 16px;">Todo App Verification</h2>
			<p>Hello,</p>
			<p>Use the following 6-digit code to finish logging in to your account:</p>
			<div style="text-align: center; margin: 28px 0;">
				<span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; background-color: #eff6ff; padding: 12px 24px; border-radius: 8px; color: #1d4ed8; border: 1px solid #bfdbfe;">%s</span>
			</div>
			<p style="color: #6b7280; font-size: 14px;">This code will expire in <strong>5 minutes</strong>. If you did not attempt to log in, please secure your account immediately.</p>
		</div>
	`, otp)

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, smtpEmail, []string{toEmail}, msg)
	if err != nil {
		log.Printf("[2FA OTP ERROR] Failed to send email via SMTP: %v", err)
		log.Printf("[2FA OTP NOTE] If using Gmail, SMTP_PASSWORD must be a 16-character Google App Password (not your OAuth client secret or standard Gmail password).")
		return err
	}

	log.Printf("[2FA OTP SUCCESS] Verification email sent to %s", toEmail)
	return nil
}
