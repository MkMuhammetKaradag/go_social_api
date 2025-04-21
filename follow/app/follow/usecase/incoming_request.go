package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type incomingRequestsUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewIncomingRequestsUseCase(sessionRepo RedisRepository, repository Repository) IncomingRequestsUseCase {
	return &incomingRequestsUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *incomingRequestsUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) ([]*domain.User, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err
	}

	users, err := u.repository.IncomingRequests(ctx, currrentUserID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
