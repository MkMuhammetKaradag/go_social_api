package bootstrap

import (
	chat "socialmedia/notification/app/chat/handler"
	chatUseCase "socialmedia/notification/app/chat/usecase"
	notification "socialmedia/notification/app/notification/handler"
	notificationUseCase "socialmedia/notification/app/notification/usecase"
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

func SetupHTTPHandlers(repo Repository, repoMongo RepositoryMongo, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {
	getNotificationsUseCase := notificationUseCase.NewGetNotificationsUseCase(repoMongo)
	getNotificationsHandler := notification.NewGetNotificationsHandler(getNotificationsUseCase)

	getUnreadNotificationsUseCase := notificationUseCase.NewGetUnreadNotificationsUseCase(repoMongo)
	getUnreadNotificationsHandler := notification.NewGetUnreadNotificationsHandler(getUnreadNotificationsUseCase)

	markNotificationUseCase := notificationUseCase.NewMarkNotificationUseCase(repoMongo)
	markNotificationHandler := notification.NewMarkNotificationHandler(markNotificationUseCase)

	deleteNotificationUseCase := notificationUseCase.NewDeleteNotificationUseCase(repoMongo)
	deleteNotificationHandler := notification.NewDeleteNotificationHandler(deleteNotificationUseCase)

	readAllNotificationsUseCase := notificationUseCase.NewReadAllNotificationsUseCase(repoMongo)
	readAllNotificationsHandler := notification.NewReadAllNotificationsHandler(readAllNotificationsUseCase)

	deleteAllNotificationsUseCase := notificationUseCase.NewDeleteAllNotificationsUseCase(repoMongo)
	deleteAllNotificationsHandler := notification.NewDeleteAllNotificationsHandler(deleteAllNotificationsUseCase)

	return map[string]interface{}{
		"getnotifications":       getNotificationsHandler,
		"getunreadnotifications": getUnreadNotificationsHandler,
		"marknotification":       markNotificationHandler,
		"deletenotification":     deleteNotificationHandler,
		"readallnotifications":   readAllNotificationsHandler,
		"deleteallnotifications": deleteAllNotificationsHandler,
	}
}
func SetupWSHandlers(repo Repository, userRedisRepo UserRedisRepository) map[string]interface{} {

	return map[string]interface{}{}
}
