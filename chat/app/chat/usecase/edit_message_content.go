package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type editMessageContentUseCase struct {
	repository Repository
}

func NewEditMessageContentUseCase(repository Repository) EditMessageContentUseCase {
	return &editMessageContentUseCase{
		repository: repository,
	}
}

func (uc *editMessageContentUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, messageID uuid.UUID, content string) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}
	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}
	err = uc.repository.UpdateMessageContent(ctx, messageID, currentUserID, content)
	if err != nil {
		return err
	}
	return nil
}
