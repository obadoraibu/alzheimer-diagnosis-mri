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

func (s *EmailSender) SendInvEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirmation Email")

	m.SetBody("text/html", fmt.Sprintf(
		"Вас пригласили в веб-приложение для анализа МРТ-снимков:<br/><br/>"+
			"<a href='http://localhost:3000/complete-invite/%s'>Регистрация</a>",
		code,
	))

	dialer := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (s *EmailSender) SendPasswordResetEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirmation Email")

	m.SetBody("text/html", fmt.Sprintf(
		"Сброс пароля:<br/><br/>"+
			"<a href='http://localhost:3000/complete-reset/%s'>Сброс</a>",
		code,
	))

	dialer := gomail.NewDialer(s.config.Host, s.config.Port, s.config.From, s.config.Password)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
