package consumer

import (
	"fmt"
	"log"
	"socialmedia/shared/messaging"
)

type EmailData struct {
	ActivationCode string
	UserName       string
}

func StartEmailConsumer(handler func(messaging.Message) error) (*messaging.RabbitMQ, error) {
	messageConfig := messaging.NewDefaultConfig()
	// messageConfig.RetryTypes = []messaging.MessageType{messaging.UserTypes.UserCreated}
	rabbit, err := messaging.NewRabbitMQ(messageConfig, messaging.EmailService)
	if err != nil {
		log.Fatal("RabbitMQ bağlantı hatası:", err)
	}

	go func() {

		err = rabbit.ConsumeMessages(func(msg messaging.Message) error {

			return handler(msg)

		})
		if err != nil {
			log.Fatal("Mesaj dinleyici başlatılamadı:", err)
		}

	}()
	return rabbit, nil
}
func handleSendEmail(msg messaging.Message) error {
	data, ok := msg.Data.(map[string]interface{})
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
	switch msg.Type {
	case "active_user":
		subject = "Hesap Aktivasyonu"
	case "forgot_password":
		subject = "Şifre Sıfırlama"
	default:
		log.Printf("Desteklenmeyen komut: %v", msg.Type)
	}

	emailData := EmailData{
		ActivationCode: activationCode,
		UserName:       userName,
	}

	log.Printf("E-posta başarıyla gönderildi. Alıcı: %s     ,  konu:%s  email,%s   templateName:%s", emailData, subject, email, templateName)
	return nil
}
