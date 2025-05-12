package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type renameConversationUseCase struct {
	repository Repository
}

func NewRenameConversationUseCase(repository Repository) RenameConversationUseCase {
	return &renameConversationUseCase{
		repository: repository,
	}
}

func (uc *renameConversationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID uuid.UUID, conversationName string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}
	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}
	err = uc.repository.UpdateConversationName(ctx, conversationID, currentUserID, conversationName)
	if err != nil {
		return err
	}
	return nil
}
