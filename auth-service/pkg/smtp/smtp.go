package smtp

import (
	"fmt"

	"github.com/go-gomail/gomail"
	"github.com/obadoraibu/go-auth/internal/config"
)

type EmailSender struct {
	config *config.SmtpConfig
}

func NewEmailSender(cfg *config.SmtpConfig) *EmailSender {
	return &EmailSender{
		config: cfg,
	}
}

func (s *EmailSender) SendConfirmationEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirmation Email")

	m.SetBody("text/html", fmt.Sprintf(
		"Thank you for signing up! Please confirm your email by clicking the following link:<br/><br/>"+
			"<a href='http://localhost:8080/email-confirm/%s'>Confirm Email</a>",
		code,
	))

	dialer := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
