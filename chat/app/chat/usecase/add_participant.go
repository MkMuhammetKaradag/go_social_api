package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type addParticipantUseCase struct {
	repository    Repository
	chatRedisRepo ChatRedisRepository
}

func NewAddParticipantUseCase(repository Repository, chatRedisRepo ChatRedisRepository) AddParticipantUseCase {
	return &addParticipantUseCase{
		repository:    repository,
		chatRedisRepo: chatRedisRepo,
	}
}

func (uc *addParticipantUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID, userID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	isblock, _ := uc.repository.HasBlockRelationship(ctx, currentUserID, userID)
	if isblock {
		return domain.ErrBlockedUser
	}
	err = uc.repository.AddParticipant(ctx, conversationID, userID, currentUserID)
	if err != nil {
		return err
	}
	user, err := uc.repository.GetUserInfoByID(ctx, userID)

	notification := &domain.ConversationUserManager{
		ConversationID: conversationID,
		UserID:         userID,
		Username:       user.Username,
		Avatar:         user.Avatar,
		Reason:         "user added in  conversation",
		Type:           "add",
	}
	err = uc.chatRedisRepo.PublishKickUserConversation(ctx, "conversation_user_manager", notification)
	if err != nil {
		fmt.Printf("Error publishing message to Redis: %v\n", err)
	}
	return nil
}
