package initializer

import (
	"log"
	"socialmedia/shared/messaging"
)

func InitMessaging() *messaging.RabbitMQ {
	config := messaging.NewDefaultConfig()
	config.RetryTypes = []messaging.MessageType{messaging.UserTypes.UserCreated}

	rabbitMQ, err := messaging.NewRabbitMQ(config, messaging.AuthService)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantısı kurulamadı: %v", err)
	}
	return rabbitMQ
}
