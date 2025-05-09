// status_hub.go
package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// UserStatusNotification, kullanıcı durum bildirimlerinin yapısı
type ConversationUserManager struct {
	UserID         uuid.UUID `json:"user_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Username       string    `json:"username"`
	Avatar         string    `json:"avatar"`
	Reason         string    `json:"reason,omitempty"`
	Type           string    `json:"type"`
}

// UserKickHub, kullanıcı durumu işlemlerini yöneten bileşen
type ConversationUserManagerHub struct {
	redisClient *redis.Client
	parentHub   *Hub
}

// NewUserKickHub, yeni bir UserKickHub örneği oluşturur
func NewConversationUserManagerHub(redisClient *redis.Client, parent *Hub) *ConversationUserManagerHub {
	return &ConversationUserManagerHub{
		redisClient: redisClient,
		parentHub:   parent,
	}
}

// Run, durum dinleme ve işleme işlevini başlatır
func (sh *ConversationUserManagerHub) Run(ctx context.Context) {

	go sh.conversationUserManager(ctx)

}

// listenForStatusUpdates, Redis'ten gelen kullanıcı durumu güncellemelerini dinler
func (sh *ConversationUserManagerHub) conversationUserManager(ctx context.Context) {
	pubsub := sh.redisClient.Subscribe(ctx, "conversation_user_manager")
	defer pubsub.Close()

	// Kanal dinleme döngüsü
	for {
		select {
		case <-ctx.Done():
			// Bağlam iptal edildiğinde çık
			return

		default:
			// Redis'ten mesaj al
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Println("Redis coversation user manager  subscription error:", err)
				continue
			}

			// Durumu işle
			var managerMsg ConversationUserManager
			err = json.Unmarshal([]byte(msg.Payload), &managerMsg)
			if err != nil {
				log.Println("coversation user manager  unmarshal error:", err)
				continue
			}

			switch managerMsg.Type {
			case "remove":
				sh.parentHub.KickUserFromConversation(ctx, managerMsg.ConversationID, managerMsg.UserID)
			case "add":
				user := UserInfo{
					Username: managerMsg.Username,
					Avatar:   managerMsg.Avatar,
					Active:   false,
				}
				sh.parentHub.AddUserToConversation(ctx, managerMsg.ConversationID, managerMsg.UserID, user)
			default:
				log.Println("Unknown conversation user manager type:", managerMsg.Type)
			}
		}
	}
}
