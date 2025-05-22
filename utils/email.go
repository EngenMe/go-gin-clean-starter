package utils

import (
	"gopkg.in/gomail.v2"

	"github.com/Caknoooo/go-gin-clean-starter/config"
)

// Dialer is an interface that matches gomail.Dialer for mocking purposes
type Dialer interface {
	DialAndSend(...*gomail.Message) error
}

// NewDialer is a variable that holds the function to create a new Dialer
var NewDialer = func(host string, port int, username, password string) Dialer {
	return gomail.NewDialer(host, port, username, password)
}

// SendMail sends an email to the specified recipient with a subject and body using pre-configured SMTP settings.
// It initializes the email configuration, constructs the email message, and sends it using a custom SMTP dialer.
// Returns an error if any step fails during the process.
func SendMail(toEmail string, subject string, body string) error {
	emailConfig, err := config.NewEmailConfig()
	if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailConfig.AuthEmail)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthEmail,
		emailConfig.AuthPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
