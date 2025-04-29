package bootstrap

import (
	"context"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/chat/domain"
	"socialmedia/chat/internal/initializer"
	"socialmedia/chat/pkg/config"

	// wsRepository "socialmedia/chat/infra/websocket"

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
	UpdateUser(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error

	CreateConversation(ctx context.Context, currrentUserID uuid.UUID, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, error)
	CreateMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs []string, attachmentTypes []string) (*domain.Message, error)
	IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)

	// IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
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
	RegisterClient(client *domain.Client, userID uuid.UUID)
	UnregisterClient(client *domain.Client, userID uuid.UUID)
	LoadConversationMembers(ctx context.Context, conversationID uuid.UUID, repo usecase.Repository) error
	SendInitialUserStatuses(client *domain.Client, conversationID uuid.UUID)
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
