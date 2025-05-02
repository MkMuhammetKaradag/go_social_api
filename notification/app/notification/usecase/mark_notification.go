package usecase

import (
	"context"
	"errors"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type markNotificationUseCase struct {
	repository Repository
}

func NewMarkNotificationUseCase(repository Repository) MarkNotificationUseCase {
	return &markNotificationUseCase{
		repository: repository,
	}
}

func (u *markNotificationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, notificationID string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return errors.New("not authorization")
	}

	userID := userData["id"]

	err := u.repository.MarkNotificationAsRead(ctx, notificationID, userID)
	if err != nil {
		return err

	}

	return nil
}
