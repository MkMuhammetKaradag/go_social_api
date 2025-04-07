package initializer

import (
	"log"
	"socialmedia/shared/messaging"
)

func InitMessaging() *messaging.RabbitMQ {
	config := messaging.NewDefaultConfig()
	config.RetryTypes = []string{"user_created"}

	rabbitMQ, err := messaging.NewRabbitMQ(config, messaging.AuthService)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantısı kurulamadı: %v", err)
	}
	return rabbitMQ
}
