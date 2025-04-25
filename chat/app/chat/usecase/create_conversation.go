package usecase

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createConversationUseCase struct {
	repository Repository
}

func NewCreateConversationUseCase(repository Repository) CreateConversationUseCase {
	return &createConversationUseCase{
		repository: repository,
	}
}

func (u *createConversationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, userIDs []uuid.UUID, name string, isGroup bool) error {

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currrentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}
	userIDs = append(userIDs, currrentUserID)

	_, err = u.repository.CreateConversation(ctx, currrentUserID, isGroup, name, userIDs)
	if err != nil {
		return err

	}

	return nil
}
