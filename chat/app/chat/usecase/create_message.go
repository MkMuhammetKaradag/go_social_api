package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type createMessageUseCase struct {
	repo          Repository
	chatRedisRepo ChatRedisRepository
}

func NewCreateMessageUseCase(repo Repository, chatRedisRepo ChatRedisRepository) CreateMessageUseCase {
	return &createMessageUseCase{repo: repo, chatRedisRepo: chatRedisRepo}
}

func (uc *createMessageUseCase) Execute(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs, attachmentTypes []string) (uuid.UUID, error) {
	// Validate input

	if conversationID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("conversation ID cannot be empty")
	}
	if senderID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("sender ID cannot be empty")
	}
	if content == "" && len(attachmentURLs) == 0 {
		return uuid.Nil, fmt.Errorf("message must contain content or attachments")
	}

	// Create the message
	message, err := uc.repo.CreateMessage(ctx, conversationID, senderID, content, attachmentURLs, attachmentTypes)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create message: %w", err)
	}
	channelName := fmt.Sprintf("conversation:%s", message.ConversationID)
	err = uc.chatRedisRepo.PublishChatMessage(ctx, channelName, message)
	if err != nil {
		fmt.Printf("Error publishing message to Redis: %v\n", err)
	}
	return message.ID, nil
}
