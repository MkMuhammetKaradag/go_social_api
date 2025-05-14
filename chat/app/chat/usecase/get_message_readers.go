package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type getMessageReadersUseCase struct {
	repository Repository
}

func NewGetMessageReadersUseCase(repository Repository) GetMessageReadersUseCase {
	return &getMessageReadersUseCase{
		repository: repository,
	}

}

func (uc *getMessageReadersUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, messageID uuid.UUID) ([]domain.User, error) {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return nil, domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return nil, err
	}

	users, err := uc.repository.GetMessageReaders(ctx, messageID, currentUserID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
