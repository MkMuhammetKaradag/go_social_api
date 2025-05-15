package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type deleteAllMessagesFromConversationUseCase struct {
	repository    Repository
	chatRedisRepo ChatRedisRepository
}

func NewDeleteAllMessagesFromConversationUseCase(repository Repository, chatRedisRepo ChatRedisRepository) DeleteAllMessagesFromConversationUseCase {
	return &deleteAllMessagesFromConversationUseCase{
		repository:    repository,
		chatRedisRepo: chatRedisRepo,
	}

}

func (uc *deleteAllMessagesFromConversationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = uc.repository.DeleteAllMessagesFromConversation(ctx, conversationID, currentUserID)
	if err != nil {
		return err
	}

	return nil
}
