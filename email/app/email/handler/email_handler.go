package handler

import (
	"bytes"

	"socialmedia/email/internal/domain"
	"socialmedia/shared/messaging"

	"fmt"
	"log"
	"text/template"
)

type EmailHandler struct {
	mailer domain.Mailer
}

type EmailData struct {
	ActivationCode string
	UserName       string
}

func NewEmailHandler(mailer domain.Mailer) *EmailHandler {
	return &EmailHandler{
		mailer: mailer,
	}
}

func (h *EmailHandler) HandleMessage(message messaging.Message) error {
	fmt.Println("handelara geldi")
	switch message.Type {
	case messaging.EmailTypes.ActivateUser:
		return h.sendActivationEmail(message)
	case messaging.EmailTypes.ForgotPassword:
		return h.sendPasswordResetEmail(message)
	// case domain.EmailTypeNotification:
	// 	return h.sendNotificationEmail(emailReq)
	default:

		return nil
	}
}

func (h *EmailHandler) sendActivationEmail(req messaging.Message) error {
	fmt.Println("sendActivationEmail  geldi")
	data, ok := req.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("geçersiz mesaj formatı")
	}

	email, emailOk := data["email"].(string)
	activationCode, codeOk := data["activation_code"].(string)
	templateName, templateOk := data["template_name"].(string)
	userName, userNameOk := data["userName"].(string)

	if !emailOk || !codeOk || !templateOk || !userNameOk {
		log.Printf("Eksik email, aktivasyon kodu veya şablon adı: %+v", data)
	}

	var subject string
	switch req.Type {
	case messaging.EmailTypes.ActivateUser:
		subject = "Hesap Aktivasyonu"
	case messaging.EmailTypes.ForgotPassword:
		subject = "Şifre Sıfırlama"
	default:
		log.Printf("Desteklenmeyen komut: %v", req.Type)
	}

	emailData := EmailData{
		ActivationCode: activationCode,
		UserName:       userName,
	}
	body, err := renderTemplate("templates/"+templateName, emailData)
	if err != nil {
		log.Printf("Şablon oluşturulamadı: %v", err)
	}
	// fmt.Println(email, subject, body)

	return h.mailer.Send(email, subject, body)
}

func (h *EmailHandler) sendPasswordResetEmail(req messaging.Message) error {
	return h.mailer.Send("req", "Hesap Aktivasyonu", "activation-template.html")
}
func renderTemplate(templatePath string, data EmailData) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
