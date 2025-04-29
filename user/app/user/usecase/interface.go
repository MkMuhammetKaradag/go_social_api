package usecase

import (
	"context"
	"socialmedia/shared/messaging"
	"socialmedia/user/domain"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProfileUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context) (*domain.User, error)
}
type CreateUserUseCase interface {
	Execute(ctx context.Context, userID, userName, email string) error
}
type GetUserUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, Identifier uuid.UUID) (*domain.User, error)
}
type SearchUserUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, Identifier string, page, limit int) ([]*domain.UserSearchResult, error)
}
type UpdateUserUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, updateuser domain.UserUpdate) error
}

type UpdateAvatarUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, avatarUrl string) error
}

type UpdateBannerUseCase interface {
	Execute(fbrCtx *fiber.Ctx, ctx context.Context, bannerUrl string) error
}

type UserStatusPublishUseCase interface {
	Execute(c *websocketFiber.Conn, ctx context.Context, currentUserID uuid.UUID)
}

type RabbitMQ interface {
	PublishMessage(ctx context.Context, msg messaging.Message) error
}
type Repository interface {
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
type Hub interface {
	Run()
	RegisterClient(client *domain.Client)
	UnregisterClient(client *domain.Client)
}
