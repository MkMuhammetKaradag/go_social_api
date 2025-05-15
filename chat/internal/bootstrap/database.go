package bootstrap

import (
	"context"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/chat/domain"
	"socialmedia/chat/internal/initializer"
	"socialmedia/chat/pkg/config"
	"socialmedia/chat/proto/userpb"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	// "google.golang.org/protobuf/runtime/protoimpl"
	// "google.golang.org/protobuf/types/known/wrapperspb"
)

//	type GetUserResponse struct {
//	    state         protoimpl.MessageState  `protogen:"open.v1"`
//	    Id            string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
//	    Username      string                  `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
//	    AvatarUrl     *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"` // nullable string
//	    unknownFields protoimpl.UnknownFields
//	    sizeCache     protoimpl.SizeCache
//	}
type UserClient interface {
	GetUserByID(ctx context.Context, userID string) (*userpb.GetUserResponse, error)
}
type Repository interface {
	CreateFollow(ctx context.Context, followerID, followingID uuid.UUID, status string) error
	DeleteFollow(ctx context.Context, followerID, followingID uuid.UUID) error
	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	HasBlockRelationship(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)

	CreateUser(ctx context.Context, id, username string) error
	UpdateUser(ctx context.Context, userID uuid.UUID, userName, avatarURL *string, isPrivate *bool) error

	CreateConversation(ctx context.Context, currrentUserID uuid.UUID, isGroup bool, name string, userIDs []uuid.UUID) (*domain.Conversation, *[]domain.BlockedParticipant, error)
	CreateMessage(ctx context.Context, conversationID, senderID uuid.UUID, content string, attachmentURLs []string, attachmentTypes []string) (*domain.Message, *domain.User, error)

	IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	GetUserIfParticipant(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) (*domain.User, error)

	// IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]domain.User, error)
	IsBlocked(ctx context.Context, userID, targetID uuid.UUID) (bool, error)
	AddParticipant(ctx context.Context, conversationID, userID, addedByUserID uuid.UUID) error
	RemoveParticipant(ctx context.Context, conversationID, userID, addedByUserID uuid.UUID) error

	PromoteToAdmin(ctx context.Context, conversationID, targetUserID, currentUserID uuid.UUID) error
	DemoteFromAdmin(ctx context.Context, conversationID, targetUserID, currentUserID uuid.UUID) error
	GetUserInfoByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)

	DeleteMessage(ctx context.Context, messageID, currentUserID uuid.UUID) (uuid.UUID, error)
	UpdateConversationName(ctx context.Context, conversationID, userID uuid.UUID, newName string) error
	UpdateMessageContent(ctx context.Context, messageID, senderID uuid.UUID, newContent string) (uuid.UUID, error)
	MarkMessagesAsRead(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) error
	MarkConversationMessagesAsRead(ctx context.Context, conversationID, userID uuid.UUID) error
	GetMessagesForConversation(ctx context.Context, conversationID, userID uuid.UUID, skip, limit int64) ([]domain.Message, error)
	GetMessageReaders(ctx context.Context, messageID, currentUserID uuid.UUID) ([]domain.User, error)
	DeleteAllMessagesFromConversation(ctx context.Context, conversationID, currentUserID uuid.UUID) error
}
type RedisRepository interface {
	GetSession(ctx context.Context, key string) (map[string]string, error)
}
type ChatRedisRepository interface {
	PublishChatMessage(ctx context.Context, channelName string, message *domain.MessageNotification) error
	GetRedisClient() *redis.Client
	PublishKickUserConversation(ctx context.Context, channelName string, message *domain.ConversationUserManager) error
}
type Hub interface {
	// Run()
	// ListenRedisSendMessage(ctx context.Context, channelName string)
	// RegisterClient(client *domain.Client, userID uuid.UUID)
	// UnregisterClient(client *domain.Client, userID uuid.UUID)
	// LoadConversationMembers(ctx context.Context, conversationID uuid.UUID, repo usecase.Repository) error
	// SendInitialUserStatuses(client *domain.Client, conversationID uuid.UUID)
	// IsConversationLoaded(conversationID uuid.UUID) bool
	// GetConversationUsers(conversationID uuid.UUID) map[string]bool

	RegisterClient(client *domain.Client, userID uuid.UUID)
	UnregisterClient(client *domain.Client, userID uuid.UUID)
	LoadConversationMembers(ctx context.Context, conversationID uuid.UUID, repo usecase.Repository) error
	IsConversationLoaded(conversationID uuid.UUID) bool
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
func InitWebsocket(ctx context.Context, redisClient *redis.Client, repo Repository) Hub {
	return initializer.InitWebsocket(ctx, redisClient, repo)
}
func InitUserClient(config config.Config) UserClient {
	return initializer.InitUserClient(config)
}
