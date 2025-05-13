package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type markMessagesAsReadUseCase struct {
	repository Repository
}

func NewMarkMessagesAsReadUseCase(repository Repository) MarkMessagesAsReadUseCase {
	return &markMessagesAsReadUseCase{
		repository: repository,
	}
}

func (uc *markMessagesAsReadUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, messagesIDs []uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = uc.repository.MarkMessagesAsRead(ctx, messagesIDs, currentUserID)
	if err != nil {
		return err
	}

	return nil
}
