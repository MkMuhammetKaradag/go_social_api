package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type promoteToAdminUseCase struct {
	repository Repository
}

func NewPromoteToAdminUseCase(repository Repository) PromoteToAdminUseCase {
	return &promoteToAdminUseCase{
		repository: repository,
	}
}

func (uc *promoteToAdminUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID, userID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = uc.repository.PromoteToAdmin(ctx, conversationID, userID, currentUserID)
	if err != nil {
		return err
	}

	return nil
}
