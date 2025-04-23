package usecase

import (
	"context"
	"socialmedia/shared/messaging"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, userID, userName string) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	CreateUser(ctx context.Context, id, username  string) error
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
