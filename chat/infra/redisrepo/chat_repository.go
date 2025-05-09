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

//	type MessageNotification struct {
//		MessageID      uuid.UUID        `json:"message_id"`
//		ConversationID uuid.UUID        `json:"conversation_id"`
//		UserID         uuid.UUID        `json:"user_id"`
//		Content        string           `json:"content"`
//		CreatedAt      string           `json:"created_at"`
//		HasAttachments bool             `json:"has_attachments"`
//		Attachments    []AttachmentInfo `json:"attachments,omitempty"`
//	}
type ConversationUserManager struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	Avatar         string    `json:"avatar"`
	Reson          string    `json:"reson"`
	Type           string    `json:"type"`
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

func (r *ChatRedisRepository) PublishChatMessage(ctx context.Context, channelName string, notification *domain.MessageNotification) error {

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
func (r *ChatRedisRepository) PublishKickUserConversation(ctx context.Context, channelName string, message *domain.ConversationUserManager) error {

	// Bildirim nesnesini oluştur
	notification := ConversationUserManager{
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		Avatar:         message.Avatar,
		Username:       message.Username,
		Reson:          message.Reason,
		Type:           message.Type,
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
