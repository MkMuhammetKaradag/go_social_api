package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type outgoingRequestsUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewOutgoingRequestsUseCase(sessionRepo RedisRepository, repository Repository) OutgoingRequestsUseCase {
	return &outgoingRequestsUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *outgoingRequestsUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) ([]*domain.FollowRequestUser, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err
	}

	users, err := u.repository.OutgoingRequests(ctx, currrentUserID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
