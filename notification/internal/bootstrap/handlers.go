package bootstrap

import (
	chat "socialmedia/notification/app/chat/handler"
	chatUseCase "socialmedia/notification/app/chat/usecase"
	user "socialmedia/notification/app/user/handler"
	userUseCase "socialmedia/notification/app/user/usecase"
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, repoMongo RepositoryMongo, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	createUserUseCase := userUseCase.NewCreateUserUseCase(repo, repoMongo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)

	chatNotificationUseCase := chatUseCase.NewChatNotificationUseCase(repoMongo)
	chatNotificationHandler := chat.NewChatNotificationHandler(chatNotificationUseCase)

	return map[messaging.MessageType]MessageHandler{
		messaging.UserTypes.UserCreated:                    createUserHandler,
		messaging.ChatTypes.UserBlockedInGroupConversation: chatNotificationHandler,
	}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {

	return map[string]interface{}{}
}
func SetupWSHandlers(repo Repository, userRedisRepo UserRedisRepository) map[string]interface{} {

	return map[string]interface{}{}
}
