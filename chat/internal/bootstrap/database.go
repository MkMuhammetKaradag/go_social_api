package bootstrap

import (
	"context"
	"socialmedia/chat/domain"
	"socialmedia/chat/internal/initializer"
	"socialmedia/chat/pkg/config"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
	CreateUser(ctx context.Context, id, username string) error

	CreateConversation(ctx context.Context, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error)
	CreateMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs []string, attachmentTypes []string) (*domain.Message, error)
}
type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
type ChatRedisRepository interface {
	PublishChatMessage(ctx context.Context, channelName string, message *domain.Message) error
	GetRedisClient() *redis.Client
}
type Hub interface {
	Run()
	ListenRedisSendMessage(ctx context.Context, channelName string)
	RegisterClient(client *domain.Client)
	UnregisterClient(client *domain.Client)
}

func InitDatabase(config config.Config) Repository {
	return initializer.InitDatabase(config)
}

func InitRedis(config config.Config) RedisRepository {
	return initializer.InitRedis(config)
}
func InitChatRedis(config config.Config) ChatRedisRepository {
	return initializer.InitChatRedis(config)
}
func InitWebsocket(redisClient *redis.Client) Hub {
	return initializer.InitWebsocket(redisClient)
}
