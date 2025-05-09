package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// MessageNotification, mesaj bildirimlerinin yapısı
type MessageNotification struct {
	Type           string           `json:"type"` // "add", "delete"
	MessageID      uuid.UUID        `json:"message_id"`
	ConversationID uuid.UUID        `json:"conversation_id"`
	UserID         uuid.UUID        `json:"user_id"`
	Username       string           `json:"username,omitempty"`
	Avatar         string           `json:"avatar,omitempty"`
	Content        string           `json:"content,omitempty"`
	CreatedAt      string           `json:"created_at,omitempty"`
	HasAttachments bool             `json:"has_attachments"`
	DeletedAt      string           `json:"deleted_at,omitempty"`
	Attachments    []AttachmentInfo `json:"attachments,omitempty"`
}

// AttachmentInfo, mesaj eklerinin yapısı
type AttachmentInfo struct {
	ID       uuid.UUID `json:"id"`
	FileURL  string    `json:"file_url"`
	FileType string    `json:"file_type"`
}

// MessageHub, mesaj işlemlerini yöneten bileşen
type MessageHub struct {
	redisClient *redis.Client
	parentHub   *Hub
}

// NewMessageHub, yeni bir MessageHub örneği oluşturur
func NewMessageHub(redisClient *redis.Client, parent *Hub) *MessageHub {
	return &MessageHub{
		redisClient: redisClient,
		parentHub:   parent,
	}
}

// Run, mesaj dinleme ve işleme işlevini başlatır
func (mh *MessageHub) Run(ctx context.Context) {
	// Redis'ten mesajları dinle
	mh.listenForMessages(ctx, "messages")
}

// listenForMessages, Redis'ten gelen mesajları dinler ve ilgili sohbetlere yönlendirir
func (mh *MessageHub) listenForMessages(ctx context.Context, channelName string) {
	pubsub := mh.redisClient.Subscribe(ctx, channelName)
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
				log.Println("Redis message subscription error:", err)
				continue
			}

			// Mesajı işle
			var notification MessageNotification
			err = json.Unmarshal([]byte(msg.Payload), &notification)
			if err != nil {
				log.Println("Message unmarshal error:", err)
				continue
			}
			fmt.Println(notification)
			// Mesaj tipini belirt
			// notification.Type = "message"

			// İlgili sohbetteki tüm istemcilere gönder
			mh.parentHub.BroadcastToConversation(ctx, notification.ConversationID, notification)
		}
	}
}

// SendMessage, bir mesajı belirli bir sohbete gönderir
func (mh *MessageHub) SendMessage(ctx context.Context, notification MessageNotification) error {
	// Mesaj tipini belirt
	notification.Type = "message"

	// Mesajı Redis'e yayınla
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return mh.redisClient.Publish(ctx, "messages", string(data)).Err()
}
