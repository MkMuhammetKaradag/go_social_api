package websocket

import "context"

type SessionRedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}

type ChatRedisRepository interface {
	PublishChatMessage(ctx context.Context, chatID string, content string, senderID string) error
}
type Repository interface {
	// CreateConversation(ctx context.Context, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error)
}
