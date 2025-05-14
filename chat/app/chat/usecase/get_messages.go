package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type getMessagesUseCase struct {
	repository Repository
}

func NewGetMessagesUseCase(repository Repository) GetMessagesUseCase {
	return &getMessagesUseCase{
		repository: repository,
	}

}

func (uc *getMessagesUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID uuid.UUID, limit, skip int64) ([]domain.Message, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err
	}

	messages, err := uc.repository.GetMessagesForConversation(ctx, conversationID, currentUserID, skip, limit)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
