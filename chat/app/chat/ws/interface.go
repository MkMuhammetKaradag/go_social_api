package ws

import (
	"context"
	"socialmedia/chat/domain"

	"github.com/google/uuid"
)

type Repository interface {
	IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]domain.User, error)
	IsBlocked(ctx context.Context, userID, targetID uuid.UUID) (bool, error)
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
