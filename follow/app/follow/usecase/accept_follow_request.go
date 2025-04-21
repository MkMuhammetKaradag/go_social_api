package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type acceptFollowRequestUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
	rabbitMQ    RabbitMQ
}

func NewAcceptFollowRequestUseCase(sessionRepo RedisRepository, repository Repository, rabbitMQ RabbitMQ) AcceptFollowRequestUseCase {
	return &acceptFollowRequestUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
		rabbitMQ:    rabbitMQ,
	}
}

func (u *acceptFollowRequestUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, requestID uuid.UUID) (string, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return "", domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return "", err
	}

	err = u.repository.AcceptFollowRequest(ctx, requestID, currrentUserID)
	if err != nil {
		return "", err
	}

	return "Follow request  accepted", nil

}
