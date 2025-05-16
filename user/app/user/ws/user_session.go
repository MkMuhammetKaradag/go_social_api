package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserLogoutNotification struct {
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Timestamp int64     `json:"timestamp"`
}

func (u UserLogoutNotification) GetUserID() uuid.UUID {
	return u.UserID
}

type SessionHub struct {
	redisClient *redis.Client
	parentHub   *Hub

	cleanupTicker *time.Ticker
}

func NewSessionHub(redisClient *redis.Client, parent *Hub) *SessionHub {
	return &SessionHub{
		redisClient:   redisClient,
		parentHub:     parent,
		cleanupTicker: time.NewTicker(30 * time.Minute), // 30 dakikada bir eskimiş durumları temizle
	}
}

// Run, durum dinleme ve işleme işlevini başlatır
func (sh *SessionHub) Run(ctx context.Context) {
	// Redis'ten kullanıcı durumlarını dinle
	go sh.listenForUserLogout(ctx)

}

func (sh *SessionHub) listenForUserLogout(ctx context.Context) {
	pubsub := sh.redisClient.Subscribe(ctx, "user:logout")
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
			var sessionNotif UserLogoutNotification
			err = json.Unmarshal([]byte(msg.Payload), &sessionNotif)
			if err != nil {
				log.Println("Status unmarshal error:", err)
				continue
			}

			// Bildirim türünü belirt
			sessionNotif.Type = "user_logout"

			fmt.Println(sessionNotif)
			sh.parentHub.BroadcastToUser(ctx, sessionNotif.UserID, sessionNotif)
		}
	}
}
