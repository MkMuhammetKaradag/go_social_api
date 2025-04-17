package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type unblockUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewUnblockUserUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) UnblockUserUseCase {
	return &unblockUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *unblockUserUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, BlockedID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = u.repository.UnblockUser(ctx, currrentUserID, BlockedID)
	if err != nil {
		return err
	}

	blockMessage := messaging.Message{
		Type:       messaging.UserTypes.UserUnBlocked,
		ToServices: []messaging.ServiceType{messaging.UserService},
		Data: map[string]interface{}{
			"unblocker_id": currrentUserID,
			"unblocked_id": BlockedID,
		},
		Critical: true,
	}

	if err := u.rabbitMQ.PublishMessage(ctx, blockMessage); err != nil {
		return err
	}
	return nil
}
