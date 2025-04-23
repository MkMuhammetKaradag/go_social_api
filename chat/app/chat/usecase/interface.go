package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/messaging"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateConversationUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, userIDs []uuid.UUID, name string, isGroup bool) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	CreateConversation(ctx context.Context, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error)
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
