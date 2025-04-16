package usecase

import (
	"context"
	"socialmedia/fallow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type blockUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewBlockUserUseCase(sessionRepo RedisRepository, repository Repository) BlockUserUseCase {
	return &blockUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
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
	return nil
}
