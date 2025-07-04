package usecase

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/shared/messaging"

	"github.com/gofiber/fiber/v2"
)

type GetNotificationsUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, limit, skip int64) ([]domain.Notification, error)
}
type GetUnreadNotificationsUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, limit, skip int64) ([]domain.Notification, error)
}

type MarkNotificationUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, notificationID string) error
}

type DeleteNotificationUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, notificationID string) error
}

type ReadAllNotificationsUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context) error
}

type DeleteAllNotificationsUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context) error
}
type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
	GetNotificationsByUserID(ctx context.Context, userID string, limit, skip int64) ([]domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error
	DeleteNotification(ctx context.Context, userID, notificationID string) error
	ReadAllNotificationsByUserID(ctx context.Context, userID string) error
	DeleteAllNotificationsByUserID(ctx context.Context, userID string) error
	GetUnreadNotificationsByUserID(ctx context.Context, userID string, limit, skip int64) ([]domain.Notification, error)
}

type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
