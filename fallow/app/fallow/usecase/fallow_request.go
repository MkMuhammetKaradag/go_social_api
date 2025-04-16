package usecase

import (
	"context"
	"socialmedia/fallow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type fallowRequestUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewFallowRequestUseCase(sessionRepo RedisRepository, repository Repository) FallowRequestUseCase {
	return &fallowRequestUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
	}
}

func (u *fallowRequestUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, FollowingID uuid.UUID) (string, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return "", domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return "", err
	}

	isPrivate, err := u.repository.IsPrivate(ctx, FollowingID)
	if err != nil {
		return "", err
	}

	if isPrivate {
		err = u.repository.CreateFollowRequest(ctx, currrentUserID, FollowingID)
		if err != nil {
			return "", err
		}
		return "Follow request sent", nil
	} else {
		err = u.repository.CreateFollow(ctx, currrentUserID, FollowingID)
		if err != nil {
			return "", err
		}
		return "User followed successfully", nil
	}
}
