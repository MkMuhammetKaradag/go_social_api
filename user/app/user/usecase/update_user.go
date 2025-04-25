package usecase

import (
	"context"
	"fmt"
	"log"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"
	"socialmedia/user/domain"

	"github.com/gofiber/fiber/v2"
)

type updateUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewUpdateUserUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) UpdateUserUseCase {
	return &updateUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *updateUserUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, updateuser domain.UserUpdate) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return fmt.Errorf("kullanıcıbilgisi  bulunamadı")
	}
	currrentUserID := userData["id"]
	err := u.repository.UpdateUser(ctx, currrentUserID, updateuser)
	if err != nil {
		return err

	}
	if updateuser.IsPrivate != nil {

		userUpdatedMessage := messaging.Message{
			Type:       messaging.UserTypes.UserUpdated,
			ToServices: []messaging.ServiceType{messaging.ChatService, messaging.FollowService},
			Data: map[string]interface{}{
				"user_id":    currrentUserID,
				"is_private": updateuser.IsPrivate,
			},
			Critical: true,
		}
		if err := u.rabbitMQ.PublishMessage(context.Background(), userUpdatedMessage); err != nil {
			log.Printf("User creation message could not be sent: %v", err)
		}
	}

	return nil
}
