package bootstrap

import (
	"context"
	"socialmedia/auth/internal/initializer"
	"socialmedia/auth/pkg/config"
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
	// messageRouter := func(msg messaging.Message) error {
	// 	handler, ok := handlers[msg.Type]
	// 	if !ok {
	// 		return nil
	// 	}
	// 	return handler.Handle(msg)
	// }

	return initializer.InitMessaging()
}
