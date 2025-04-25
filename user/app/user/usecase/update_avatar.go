package usecase

import (
	"context"
	"fmt"
	"log"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type updateAvatarUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewUpdateAvatarUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) UpdateAvatarUseCase {
	return &updateAvatarUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *updateAvatarUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, avatarURL string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	// userID := userData["id"]
	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}
	err = u.repository.UpdateAvatar(ctx, currrentUserID, avatarURL)
	if err != nil {
		return err

	}

	userUpdatedMessage := messaging.Message{
		Type:       messaging.UserTypes.UserUpdated,
		ToServices: []messaging.ServiceType{messaging.ChatService, messaging.FollowService},
		Data: map[string]interface{}{
			"user_id":    currrentUserID,
			"avatar_url": avatarURL,
		},
		Critical: true,
	}
	if err := u.rabbitMQ.PublishMessage(context.Background(), userUpdatedMessage); err != nil {
		log.Printf("User creation message could not be sent: %v", err)
	}

	return nil
}
