package misc

import (
	"net/smtp"
	"strings"
)

type Email struct {
	SmtpHost string
	SmtpPort string
	Username string
	Password string
}

func NewEmail(smtpHost string, smtpPort string, username string, password string) *Email {
	return &Email{
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
		Username: username,
		Password: password,
	}
}

func (e *Email) SendHTML(from string, to []string, subject, body string) error {
	server := e.SmtpHost + ":" + e.SmtpPort

	fromStr := "From: " + from + "\n"
	toStr := "To: " + strings.Join(to, ",") + "\n"
	subjectStr := "Subject: " + subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(fromStr + toStr + subjectStr + mime + "\n" + body)

	return smtp.SendMail(server, smtp.PlainAuth("", e.Username, e.Password, e.SmtpHost), from, to, msg)
}
