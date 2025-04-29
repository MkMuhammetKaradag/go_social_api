package usecase

import (
	"context"
	"fmt"
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
	fmt.Println("hello")

	userIDMap := make(map[uuid.UUID]struct{})
	for _, id := range userIDs {
		userIDMap[id] = struct{}{}
	}
	if _, exists := userIDMap[currrentUserID]; !exists {
		userIDs = append(userIDs, currrentUserID)
	}
	fmt.Println("exists")
	if len(userIDs) <= 2 {
		isGroup = false
	}
	_, err = u.repository.CreateConversation(ctx, currrentUserID, isGroup, name, userIDs)
	if err != nil {
		return err

	}

	return nil
}
