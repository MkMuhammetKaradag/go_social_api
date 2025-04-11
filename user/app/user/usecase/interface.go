package usecase

import (
	"context"
	"socialmedia/shared/messaging"

	"github.com/gofiber/fiber/v2"
)

type ProfileUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
