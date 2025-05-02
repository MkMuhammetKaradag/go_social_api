package usecase

import (
	"context"
	"errors"
	"socialmedia/notification/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
)

type getNotificationsUseCase struct {
	repository Repository
}

func NewGetNotificationsUseCase(repository Repository) GetNotificationsUseCase {
	return &getNotificationsUseCase{
		repository: repository,
	}
}

func (u *getNotificationsUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, limit, skip int64) ([]domain.Notification, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, errors.New("not authorization")
	}

	userID := userData["id"]

	notification, err := u.repository.GetNotificationsByUserID(ctx, userID, limit, skip)
	if err != nil {
		return nil, err

	}

	return notification, nil
}
