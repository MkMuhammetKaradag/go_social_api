package usecase

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type readAllNotificationsUseCase struct {
	repository Repository
}

func NewReadAllNotificationsUseCase(repository Repository) ReadAllNotificationsUseCase {
	return &readAllNotificationsUseCase{
		repository: repository,
	}
}

func (u *readAllNotificationsUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}
	userID := userData["id"]

	err := u.repository.ReadAllNotificationsByUserID(ctx, userID)
	if err != nil {
		return nil
	}
	return nil
}
