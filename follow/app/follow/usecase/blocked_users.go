package usecase

import (
	"context"
	"socialmedia/follow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type getBlockedUsersUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewGetBlockedUsersUseCase(sessionRepo RedisRepository, repository Repository) GetBlockedUsersUseCase {
	return &getBlockedUsersUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *getBlockedUsersUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context) ([]*domain.BlockedUser, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err
	}

	users, err := u.repository.GetBlockedUsers(ctx, currrentUserID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
