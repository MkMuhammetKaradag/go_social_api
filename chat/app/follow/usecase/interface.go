package usecase

import (
	"context"
	"socialmedia/shared/messaging"

	"github.com/google/uuid"
)

type FollowRequestUseCase interface {
	Execute(ctx context.Context, followerID, followingID uuid.UUID, status string) error
}
type UnFollowRequestUseCase interface {
	Execute(ctx context.Context, followerID, followingID uuid.UUID) error
}
type BlockUserUseCase interface {
	Execute(ctx context.Context, blockerID, blockedID uuid.UUID) error
}
type UnBlockUserUseCase interface {
	Execute(ctx context.Context, blockerID, blockedID uuid.UUID) error
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
