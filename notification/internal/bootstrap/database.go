package bootstrap

import (
	"context"
	"socialmedia/notification/internal/initializer"
	"socialmedia/notification/pkg/config"

	"github.com/google/uuid"
)

type Repository interface {
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


type UserRedisRepository interface {
	PublishUserStatus(ctx context.Context, userID uuid.UUID, status string)
}
