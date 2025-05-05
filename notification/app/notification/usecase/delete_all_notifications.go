package usecase

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type deleteAllNotificationsUseCase struct {
	repository Repository
}

func NewDeleteAllNotificationsUseCase(repository Repository) DeleteAllNotificationsUseCase {
	return &deleteAllNotificationsUseCase{
		repository: repository,
	}
}

func (u *deleteAllNotificationsUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	userID := userData["id"]
	err := u.repository.DeleteAllNotificationsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
