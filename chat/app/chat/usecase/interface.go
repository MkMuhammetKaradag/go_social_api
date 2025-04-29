package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/messaging"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateConversationUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, userIDs []uuid.UUID, name string, isGroup bool) error
}
type CreateMessageUseCase interface {
	Execute(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs, attachmentTypes []string) (uuid.UUID, error)
}
type ChatWebSocketListenUseCase interface {
	Execute(c *websocketFiber.Conn, ctx context.Context, userID, conversationID uuid.UUID)
}

type Hub interface {
	Run()
	ListenRedisSendMessage(ctx context.Context, channelName string)
	RegisterClient(client *domain.Client, userID uuid.UUID)
	UnregisterClient(client *domain.Client, userID uuid.UUID)
	LoadConversationMembers(ctx context.Context, conversationID uuid.UUID, repo Repository) error
	SendInitialUserStatuses(client *domain.Client, conversationID uuid.UUID)
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}

type Repository interface {
	CreateConversation(ctx context.Context, currrentUserID uuid.UUID, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error)
	CreateMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs []string, attachmentTypes []string) (*domain.Message, error)
	IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
type ChatRedisRepository interface {
	PublishChatMessage(ctx context.Context, channelName string, message *domain.Message) error
}
