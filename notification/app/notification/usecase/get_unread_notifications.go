package usecase

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type getUnreadNotificationsUseCase struct {
	repository Repository
}

func NewGetUnreadNotificationsUseCase(repository Repository) GetUnreadNotificationsUseCase {
	return &getUnreadNotificationsUseCase{
		repository: repository,
	}
}

func (u *getUnreadNotificationsUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, limit, skip int64) ([]domain.Notification, error) {

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, domain.ErrNotFoundAuthorization
	}
	userID := userData["id"]

	notifications, err := u.repository.GetUnreadNotificationsByUserID(ctx, userID, limit, skip)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
