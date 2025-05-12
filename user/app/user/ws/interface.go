package websocket

import (
	"context"

	"github.com/google/uuid"
)

type UserRedisRepository interface {
	PublishUserStatus(ctx context.Context, userID uuid.UUID, status string)
	// GetRedisClient() *redis.Client
}
