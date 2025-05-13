package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type editMessageContentUseCase struct {
	repository    Repository
	chatRedisRepo ChatRedisRepository
}

func NewEditMessageContentUseCase(repository Repository, chatRedisRepo ChatRedisRepository) EditMessageContentUseCase {
	return &editMessageContentUseCase{
		repository:    repository,
		chatRedisRepo: chatRedisRepo,
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
	conversationID, err := uc.repository.UpdateMessageContent(ctx, messageID, currentUserID, content)
	if err != nil {
		return err
	}

	
	notification := &domain.MessageNotification{
		Type:           "message_edit",
		MessageID:      messageID,
		ConversationID: conversationID,
		UserID:         currentUserID,
		Content:        content,
	}


	err = uc.chatRedisRepo.PublishChatMessage(ctx, "messages", notification)
	if err != nil {
		fmt.Printf("Error publishing message to Redis: %v\n", err)
	}
	return nil
}
