package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type removeParticipantUseCase struct {
	repository    Repository
	chatRedisRepo ChatRedisRepository
}

func NewRemoveParticipantUseCase(repository Repository, chatRedisRepo ChatRedisRepository) RemoveParticipantUseCase {
	return &removeParticipantUseCase{
		repository:    repository,
		chatRedisRepo: chatRedisRepo,
	}

}

func (uc *removeParticipantUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, conversationID, userID uuid.UUID) error {
	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	err = uc.repository.RemoveParticipant(ctx, conversationID, userID, currentUserID)
	if err != nil {
		return err
	}
	notification := &domain.KickUserConservation{
		ConversationID: conversationID,
		UserID:         userID,
	}
	err = uc.chatRedisRepo.PublishKickUserConversation(ctx, "kick_user_channel", notification)
	if err != nil {
		fmt.Printf("Error publishing message to Redis: %v\n", err)
	}
	return nil
}
