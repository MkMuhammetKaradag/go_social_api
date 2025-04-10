package mailer

import (
	"fmt"
	"net/smtp"
	"socialmedia/email/pkg/config"
)

type SMTPMailer struct {
	host     string
	port     string
	email    string
	password string
}

func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.SMTP.Host,
		port:     cfg.SMTP.Port,
		email:    cfg.SMTP.Email,
		password: cfg.SMTP.Password,
	}
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	msg := []byte("Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", m.email, m.password, m.host)
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	fmt.Println(addr)
	err := smtp.SendMail(addr, auth, m.email, []string{to}, msg)
	if err != nil {
		fmt.Println("error aldım gönderirken", err)
		return err
	}
	return nil

}
