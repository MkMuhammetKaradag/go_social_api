package usecase

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/shared/messaging"
)

type ChatNotificationUseCase interface {
	Execute(ctx context.Context, notification domain.Notification) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}

//	type Repository interface {
//		CreateUser(ctx context.Context, userID uuid.UUID, username string) error
//	}
type Repository interface {
	CreateNotification(ctx context.Context, notification domain.Notification) error
}
type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
