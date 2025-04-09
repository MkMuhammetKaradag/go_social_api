package usecase

import (
	"context"
	"socialmedia/shared/messaging"
)

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
