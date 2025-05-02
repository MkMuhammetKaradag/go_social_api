package usecase

import (
	"context"
	"socialmedia/notification/domain"
)

type chatNotificationUseCase struct {
	repositoryMongo Repository
}

func NewChatNotificationUseCase(repositoryMongo Repository) ChatNotificationUseCase {
	return &chatNotificationUseCase{
		repositoryMongo: repositoryMongo,
	}
}

func (u *chatNotificationUseCase) Execute(ctx context.Context, notification domain.Notification) error {

	err := u.repositoryMongo.CreateNotification(ctx, notification)
	if err != nil {
		return err

	}

	return nil
}
