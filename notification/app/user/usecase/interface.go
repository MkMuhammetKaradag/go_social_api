package usecase

import (
	"context"
	"socialmedia/shared/messaging"

	"github.com/google/uuid"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, userName string) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	CreateUser(ctx context.Context, userID uuid.UUID, username string) error
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
