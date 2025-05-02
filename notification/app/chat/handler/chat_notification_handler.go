package chat

import (
	"context"
	"fmt"
	"log"
	"socialmedia/notification/app/chat/usecase"
	"socialmedia/notification/domain"
	"socialmedia/shared/messaging"
)

type ChatNotificationHandler struct {
	usecase usecase.ChatNotificationUseCase
}

func NewChatNotificationHandler(ChatNotificationUsecase usecase.ChatNotificationUseCase) *ChatNotificationHandler {
	return &ChatNotificationHandler{
		usecase: ChatNotificationUsecase,
	}
}

func (h *ChatNotificationHandler) Handle(msg messaging.Message) error {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return domain.ErrInvalidMessageFormat
	}

	actorStr, ok := data["actor_id"].(string)
	if !ok {
		return domain.ErrMissingActorID
	}

	convStr, ok := data["conversation_id"].(string)
	if !ok {
		return domain.ErrMissingConversationID
	}

	groupName, _ := data["group_name"].(string)

	rawPairs, ok := data["blocked_pairs"].([]interface{})
	if !ok {
		return domain.ErrInvalidBlockedPairs
	}

	uniqueBlockers := make(map[string]int)

	for _, item := range rawPairs {
		pairMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		blockerStr, ok1 := pairMap["blocker_id"].(string)
		if !ok1 {
			continue
		}
		uniqueBlockers[blockerStr]++
	}

	for blockerStr, i := range uniqueBlockers {

		var content string
		if i == 1 {
			content = fmt.Sprintf("Dahil olduğunuz '%s' adlı sohbette blokladığınız kullanıcı bulunmaktadır, mesajlarını göremeyeceksiniz.", groupName)
		} else {
			content = fmt.Sprintf("Dahil olduğunuz '%s' adlı sohbette blokladığınız kullanıcılar bulunmaktadır, mesajlarını göremeyeceksiniz.", groupName)
		}
		url := fmt.Sprintf("/conversation/'%s'", convStr)

		notification := domain.Notification{
			UserID:     blockerStr,
			ActorID:    actorStr,
			EntityID:   convStr,
			Type:       "block",
			EntityType: "conversation",
			Url:        url,
			GroupName:  groupName,
			Content:    content,
		}
		err := h.usecase.Execute(context.Background(), notification)
		if err != nil {
			log.Printf("notification error: %v", notification)
		}
	}
	return nil
}
