package consumer

import (
	"log"
	"socialmedia/shared/messaging"
)

func StartUserConsumer(handler func(messaging.Message) error) (*messaging.RabbitMQ, error) {
	messageConfig := messaging.NewDefaultConfig()

	rabbit, err := messaging.NewRabbitMQ(messageConfig, messaging.FollowService)
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
