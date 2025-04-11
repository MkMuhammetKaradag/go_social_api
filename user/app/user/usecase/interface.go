package usecase

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type ProfileUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context) (*domain.User, error)
}
type CreateUserUseCase interface {
	Execute(ctx context.Context, userID, userName, email string) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	GetUserProfile(ctx context.Context, identifier string) (*domain.User, error)
	CreateUser(ctx context.Context, id, username, email string) error
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
