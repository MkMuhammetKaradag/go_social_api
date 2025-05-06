package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type demoteFromAdminUseCase struct {
	repository Repository
}

func NewDemoteFromAdminUseCase(repository Repository) DemoteFromAdminUseCase {
	return &demoteFromAdminUseCase{
		repository: repository,
	}
}

func (uc *demoteFromAdminUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID, userID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = uc.repository.DemoteFromAdmin(ctx, conversationID, userID, currentUserID)
	if err != nil {
		return err
	}

	return nil
}
