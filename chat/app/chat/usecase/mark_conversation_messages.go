package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type markConversationMessagesAsReadUseCase struct {
	repository Repository
}

func NewMarkConversationMessagesAsReadUseCase(repository Repository) MarkConversationMessagesAsReadUseCase {
	return &markConversationMessagesAsReadUseCase{
		repository: repository,
	}
}

func (uc *markConversationMessagesAsReadUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = uc.repository.MarkConversationMessagesAsRead(ctx, conversationID, currentUserID)
	if err != nil {
		return err
	}

	return nil
}
