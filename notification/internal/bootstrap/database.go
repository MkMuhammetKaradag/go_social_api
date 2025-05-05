package bootstrap

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/notification/internal/initializer"
	"socialmedia/notification/pkg/config"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, userID uuid.UUID, username string) error
}
type RepositoryMongo interface {
	CreateUser(ctx context.Context, userID uuid.UUID, username string) error
	CreateNotification(ctx context.Context, notification domain.Notification) error
	GetNotificationsByUserID(ctx context.Context, userID string, limit, skip int64) ([]domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error
	DeleteNotification(ctx context.Context, userID, notificationID string) error
	ReadAllNotificationsByUserID(ctx context.Context, userID string) error
	DeleteAllNotificationsByUserID(ctx context.Context, userID string) error
}
type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}

func InitDatabase(config config.Config) Repository {
	return initializer.InitDatabase(config)
}
func InitDatabaseMongo(config config.Config) RepositoryMongo {
	return initializer.InitDatabaseMongo(config)
}

func InitRedis(config config.Config) RedisRepository {
	return initializer.InitRedis(config)
}

type UserRedisRepository interface {
	PublishUserStatus(ctx context.Context, userID uuid.UUID, status string)
}
