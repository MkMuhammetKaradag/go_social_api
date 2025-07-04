package bootstrap

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/follow/internal/initializer"
	"socialmedia/follow/pkg/config"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, id, username string) error
	UpdateUser(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error
	
	IsPrivate(ctx context.Context, userID uuid.UUID) (bool, error)
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	CreateFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	DeleteFollowRequest(ctx context.Context, requesterID, targetID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	IncomingRequests(ctx context.Context, currentUserID uuid.UUID) ([]*domain.User, error)
	OutgoingRequests(ctx context.Context, currentUserID uuid.UUID) ([]*domain.FollowRequestUser, error)
	GetBlockedUsers(ctx context.Context, currentUserID uuid.UUID) ([]*domain.BlockedUser, error)
	AcceptFollowRequest(ctx context.Context, requestID, currentUserID uuid.UUID) (uuid.UUID, error)
	RejectFollowRequest(ctx context.Context, requestID, currentUserID uuid.UUID) (uuid.UUID, error)
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
