package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type EmailService struct {
	From     string
	Password string
	Host     string
	Port     string
}

type EmailData struct {
	To       string
	Subject  string
	Body     string
	Template string
	Data     interface{}
}

func NewEmailService() *EmailService {
	return &EmailService{
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}
}

func (e *EmailService) SendEmail(data EmailData) error {
	// Load template
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/%s.html", data.Template))
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data.Data); err != nil {
		return err
	}

	// Send email
	auth := smtp.PlainAuth("", e.From, e.Password, e.Host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s",
		data.To, data.Subject, body.String()))

	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)
	return smtp.SendMail(addr, auth, e.From, []string{data.To}, msg)
}
