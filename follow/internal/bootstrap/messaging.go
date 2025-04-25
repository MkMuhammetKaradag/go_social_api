package bootstrap

import (
	"context"
	"fmt"
	"socialmedia/follow/internal/initializer"
	"socialmedia/follow/pkg/config"
	"socialmedia/shared/messaging"
)

type Messaging interface {
	Close() error
	PublishMessage(ctx context.Context, msg messaging.Message) error
}

type MessageHandler interface {
	Handle(msg messaging.Message) error
}

func SetupMessaging(handlers map[messaging.MessageType]MessageHandler, config config.Config) Messaging {
	messageRouter := func(msg messaging.Message) error {
		fmt.Println(msg.Type)
		handler, ok := handlers[msg.Type]
		if !ok {
			return nil
		}
		return handler.Handle(msg)
	}

	return initializer.InitMessaging(messageRouter)
}
