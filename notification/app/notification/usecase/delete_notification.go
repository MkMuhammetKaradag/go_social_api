package usecase

import (
	"context"
	"socialmedia/notification/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type deleteNotificationUseCase struct {
	repository Repository
}

func NewDeleteNotificationUseCase(repository Repository) DeleteNotificationUseCase {
	return &deleteNotificationUseCase{
		repository: repository,
	}
}

func (u *deleteNotificationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, notificationID string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	userID := userData["id"]
	err := u.repository.DeleteNotification(ctx, userID, notificationID)
	if err != nil {
		return err
	}

	return nil
}
