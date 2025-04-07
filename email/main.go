package main

import (
	"fmt"
	"log"
	"os"
	"socialmedia/auth/pkg/graceful"
	"socialmedia/email/internal/server"
	"socialmedia/email/pkg/config"
	"socialmedia/shared/messaging"
	"time"
)

type EmailData struct {
	ActivationCode string
	UserName       string
}

func main() {
	appConfig := config.Read()
	messageConfig := messaging.NewDefaultConfig()

	rabbit, err := messaging.NewRabbitMQ(messageConfig, messaging.EmailService)
	if err != nil {
		log.Fatal("RabbitMQ bağlantı hatası:", err)
	}
	defer rabbit.Close()

	err = rabbit.ConsumeMessages(func(msg messaging.Message) error {
		if msg.Type == "active_user" || msg.Type == "forgot_password" {
			fmt.Println(msg)
			// return nil
			return handleSendEmail(msg)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Mesaj dinleyici başlatılamadı:", err)
	}

	serverConfig := server.Config{
		Port:         appConfig.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app := server.NewFiberApp(serverConfig)
	go func() {
		if err := server.Start(app, appConfig.Server.Port); err != nil {

			os.Exit(1)
		}
	}()

	graceful.WaitForShutdown(app, 5*time.Second)

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
