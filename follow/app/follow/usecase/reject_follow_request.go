package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type rejectFollowRequestUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewRejectFollowRequestUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) RejectFollowRequestUseCase {
	return &rejectFollowRequestUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *rejectFollowRequestUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, requestID uuid.UUID) (string, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return "", domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return "", err
	}

	err = u.repository.RejectFollowRequest(ctx, requestID, currrentUserID)
	if err != nil {
		return "", err
	}

	return "Follow request  rejected", nil

}
