package usecase

import (
	"context"
	"socialmedia/shared/messaging"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FallowRequestUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, userID uuid.UUID) (string, error)
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	IsPrivate(ctx context.Context, userID uuid.UUID) (bool, error)
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	CreateFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error
	// UserExists(ctx context.Context, userID uuid.UUID) (bool, error)
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
