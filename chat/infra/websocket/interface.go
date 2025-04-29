package websocket

import (
	"context"
)

type SessionRedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}

type ChatRedisRepository interface {
	PublishChatMessage(ctx context.Context, chatID string, content string, senderID string) error
}
// type Repository interface {
// 	IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
// 	GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
// }

