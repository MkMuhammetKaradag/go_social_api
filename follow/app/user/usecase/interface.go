package usecase

import (
	"context"
	"socialmedia/shared/messaging"

	"github.com/google/uuid"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, userID, userName string) error
}
type UpdateUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	CreateUser(ctx context.Context, id, username string) error
	UpdateUser(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
