package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type blockUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewBlockUserUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) BlockUserUseCase {
	return &blockUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *blockUserUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, BlockedID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = u.repository.BlockUser(ctx, currrentUserID, BlockedID)
	if err != nil {
		return err
	}

	blockMessage := messaging.Message{
		Type:       messaging.UserTypes.UserBlocked,
		ToServices: []messaging.ServiceType{messaging.UserService, messaging.ChatService},
		Data: map[string]interface{}{
			"blocker_id": currrentUserID,
			"blocked_id": BlockedID,
		},
		Critical: true,
	}

	if err := u.rabbitMQ.PublishMessage(ctx, blockMessage); err != nil {
		return err
	}
	return nil
}
