package initializer

import (
	"log"
	"socialmedia/shared/messaging"
)

func InitMessaging(handler func(messaging.Message) error) *messaging.RabbitMQ {
	config := messaging.NewDefaultConfig()
	config.RetryTypes = []messaging.MessageType{messaging.UserTypes.UserCreated}

	rabbitMQ, err := messaging.NewRabbitMQ(config, messaging.FallowService)
	if err != nil {
		log.Fatalf("RabbitMQ bağlantısı kurulamadı: %v", err)
	}

	go func() {

		err = rabbitMQ.ConsumeMessages(func(msg messaging.Message) error {

			return handler(msg)

		})
		if err != nil {
			log.Fatal("Mesaj dinleyici başlatılamadı:", err)
		}

	}()
	return rabbitMQ
}
