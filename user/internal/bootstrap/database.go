package bootstrap

import (
	"context"
	"socialmedia/user/domain"
	"socialmedia/user/internal/initializer"
	"socialmedia/user/pkg/config"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
	GetUserProfile(ctx context.Context, identifier string) (*domain.User, error)
	CreateUser(ctx context.Context, id, username, email string) error
	UpdateUser(ctx context.Context, userID string, update domain.UserUpdate) error
	GetUser(ctx context.Context, currrentUserID, targetUserID uuid.UUID) (*domain.User, error)
	SearchUsers(ctx context.Context, currentUserID uuid.UUID, searchTerm string, page, limit int) ([]*domain.UserSearchResult, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) error
	UpdateBanner(ctx context.Context, userID uuid.UUID, bannerURL string) error
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

func InitUserRedis(config config.Config) UserRedisRepository {
	return initializer.InitUserRedis(config)
}
func InitWebsocket(ctx context.Context, redisRepo UserRedisRepository) Hub {
	client := redisRepo.GetRedisClient()
	return initializer.InitWebsocket(ctx, redisRepo, client)
}

type UserRedisRepository interface {
	PublishUserStatus(ctx context.Context, userID uuid.UUID, status string)

	GetRedisClient() *redis.Client
}
type Hub interface {
	Run(ctx context.Context)
	RegisterClient(client *domain.Client)
	UnregisterClient(client *domain.Client)
}
