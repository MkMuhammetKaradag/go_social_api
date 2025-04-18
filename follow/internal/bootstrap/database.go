package bootstrap

import (
	"context"
	"socialmedia/follow/internal/initializer"
	"socialmedia/follow/pkg/config"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, id, username string) error
	IsPrivate(ctx context.Context, userID uuid.UUID) (bool, error)
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	CreateFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	DeleteFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
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
