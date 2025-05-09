package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type deleteMessageUseCase struct {
	repository    Repository
	chatRedisRepo ChatRedisRepository
}

func NewDeleteMessageUseCase(repository Repository, chatRedisRepo ChatRedisRepository) DeleteMessageUseCase {
	return &deleteMessageUseCase{
		repository:    repository,
		chatRedisRepo: chatRedisRepo,
	}

}

func (uc *deleteMessageUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, messageID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	conversationID, err := uc.repository.DeleteMessage(ctx, messageID, currentUserID)
	if err != nil {
		return err
	}
	fmt.Println(conversationID)
	// notification := &domain.ConversationUserManager{
	// 	ConversationID: conversationID,
	// 	UserID:         userID,
	// 	Reason:         "user removed from conversation",
	// 	Type:           "remove",
	// }
	// err = uc.chatRedisRepo.PublishKickUserConversation(ctx, "conversation_user_manager", notification)
	// if err != nil {
	// 	fmt.Printf("Error publishing message to Redis: %v\n", err)
	// }
	return nil
}
