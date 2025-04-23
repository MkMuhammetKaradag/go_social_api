package bootstrap

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/chat/internal/initializer"
	"socialmedia/chat/pkg/config"

	"github.com/google/uuid"
)

type Repository interface {
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
	CreateUser(ctx context.Context, id, username string) error
	CreateConversation(ctx context.Context, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error)
}
type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}

func InitDatabase(config config.Config) Repository {
	return initializer.InitDatabase(config)
}

func InitRedis(config config.Config) RedisRepository {
	return initializer.InitRedis(config)
}
