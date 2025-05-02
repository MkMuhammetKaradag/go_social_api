package usecase

import (
	"context"
	"log"
	"socialmedia/chat/domain"
	"socialmedia/shared/messaging"
	"socialmedia/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createConversationUseCase struct {
	repository Repository
	rabbitMQ   RabbitMQ
}

func NewCreateConversationUseCase(repository Repository, rabbitMQ RabbitMQ) CreateConversationUseCase {
	return &createConversationUseCase{
		repository: repository,
		rabbitMQ:   rabbitMQ,
	}
}

func (u *createConversationUseCase) Execute(fbrCtx *fiber.Ctx, ctx context.Context, userIDs []uuid.UUID, name string, isGroup bool) error {

	userData, ok := middlewares.GetUserData(fbrCtx)
	if !ok {
		return domain.ErrNotFoundAuthorization
	}

	currentUserID, err := uuid.Parse(userData["id"])
	if err != nil {
		return err
	}

	// UUID'leri benzersiz hale getir
	uniqueUserIDMap := make(map[uuid.UUID]struct{})
	for _, id := range userIDs {
		uniqueUserIDMap[id] = struct{}{}
	}
	uniqueUserIDMap[currentUserID] = struct{}{} // current user'ı da ekle

	// Benzersiz ID listesini oluştur
	uniqueUserIDs := make([]uuid.UUID, 0, len(uniqueUserIDMap))
	for id := range uniqueUserIDMap {
		isblock, _ := u.repository.HasBlockRelationship(ctx, currentUserID, id)
		if isblock {
			continue
		}
		uniqueUserIDs = append(uniqueUserIDs, id)
	}

	// 2 veya daha az katılımcı varsa grup değildir
	isGroup = len(uniqueUserIDs) > 2

	conversation, blockedParticipant, err := u.repository.CreateConversation(ctx, currentUserID, isGroup, name, uniqueUserIDs)
	if err != nil {
		return err

	}

	blockedPairs := make([]map[string]string, 0)
	for _, b := range *blockedParticipant {
		blockedPairs = append(blockedPairs, map[string]string{
			"blocker_id": b.BlockerID.String(),
			"blocked_id": b.BlockedID.String(),
		})
	}
	if blockedParticipant != nil && len(*blockedParticipant) > 0 {
		blockedEvent := messaging.Message{
			Type:       messaging.ChatTypes.UserBlockedInGroupConversation,
			ToServices: []messaging.ServiceType{messaging.NotificationService},
			Data: map[string]interface{}{
				"entity_type":     "conversation",
				"actor_id":        currentUserID.String(),
				"conversation_id": conversation.ID.String(),
				"group_name":      name,
				"blocked_pairs":   blockedPairs,
			},
			Critical: false,
		}
		if err := u.rabbitMQ.PublishMessage(ctx, blockedEvent); err != nil {
			log.Printf("notification  message could not be sent: %v", err)
			// return err
		}
	}

	return nil
}
