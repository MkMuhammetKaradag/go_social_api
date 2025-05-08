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
type UserKickChannel struct {
	UserID         uuid.UUID `json:"user_id"`
	ConversationID uuid.UUID `json:"conversation_id"` // "online" veya "offline"
}

// UserKickHub, kullanıcı durumu işlemlerini yöneten bileşen
type UserKickHub struct {
	redisClient *redis.Client
	parentHub   *Hub
}

// NewUserKickHub, yeni bir UserKickHub örneği oluşturur
func NewUserKickHub(redisClient *redis.Client, parent *Hub) *UserKickHub {
	return &UserKickHub{
		redisClient: redisClient,
		parentHub:   parent,
	}
}

// Run, durum dinleme ve işleme işlevini başlatır
func (sh *UserKickHub) Run(ctx context.Context) {
	// Redis'ten kullanıcı durumlarını dinle
	go sh.listenForStatusUpdates(ctx)

}

// listenForStatusUpdates, Redis'ten gelen kullanıcı durumu güncellemelerini dinler
func (sh *UserKickHub) listenForStatusUpdates(ctx context.Context) {
	pubsub := sh.redisClient.Subscribe(ctx, "kick_user_channel")
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
				log.Println("Redis status subscription error:", err)
				continue
			}

			// Durumu işle
			var kickUser UserKickChannel
			err = json.Unmarshal([]byte(msg.Payload), &kickUser)
			if err != nil {
				log.Println("Status unmarshal error:", err)
				continue
			}
			sh.parentHub.KickUserFromConversation(ctx, kickUser.ConversationID, kickUser.UserID)

		}
	}
}
