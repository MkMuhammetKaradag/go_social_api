package usecase

import (
	"context"
	"socialmedia/fallow/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type unblockUserUseCase struct {
	sessionRepo RedisRepository
	repository  Repository
}

func NewUnblockUserUseCase(sessionRepo RedisRepository, repository Repository) UnblockUserUseCase {
	return &unblockUserUseCase{
		sessionRepo: sessionRepo,
		repository:  repository,
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
	return nil
}
