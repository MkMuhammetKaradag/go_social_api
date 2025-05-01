package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"socialmedia/chat/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ChatRedisRepository struct {
	client *redis.Client
}

type MessageNotification struct {
	MessageID      uuid.UUID        `json:"message_id"`
	ConversationID uuid.UUID        `json:"conversation_id"`
	UserID         uuid.UUID        `json:"user_id"`
	Content        string           `json:"content"`
	CreatedAt      string           `json:"created_at"`
	HasAttachments bool             `json:"has_attachments"`
	Attachments    []AttachmentInfo `json:"attachments,omitempty"`
}

type AttachmentInfo struct {
	ID       uuid.UUID `json:"id"`
	FileURL  string    `json:"file_url"`
	FileType string    `json:"file_type"`
}

func NewChatRedisRepository(connString, password string, db int) (*ChatRedisRepository, error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &ChatRedisRepository{client: RedisClient}, nil
}

func (r *ChatRedisRepository) PublishChatMessage(ctx context.Context, channelName string, message *domain.Message) error {
	fmt.Println("geldi")

	attachments := []AttachmentInfo{}
	for _, attachment := range message.Attachments {
		attachments = append(attachments, AttachmentInfo{
			ID:       attachment.ID,
			FileURL:  attachment.FileURL,
			FileType: attachment.FileType,
		})
	}

	// Bildirim nesnesini oluştur
	notification := MessageNotification{
		MessageID:      message.ID,
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		Content:        message.Content,
		CreatedAt:      message.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		HasAttachments: len(message.Attachments) > 0,
		Attachments:    attachments,
	}

	// JSON'a dönüştür
	notificationJson, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal message notification: %w", err)
	}

	// Redis'e yayınla
	err = r.client.Publish(ctx, channelName, notificationJson).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to Redis: %w", err)
	}
	return nil
}
func (r *ChatRedisRepository) GetRedisClient() *redis.Client {
	return r.client
}
