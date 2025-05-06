package bootstrap

import (
	chat "socialmedia/chat/app/chat/handler"
	chatUseCase "socialmedia/chat/app/chat/usecase"
	follow "socialmedia/chat/app/follow/handler"
	followUseCase "socialmedia/chat/app/follow/usecase"
	user "socialmedia/chat/app/user/handler"
	userUseCase "socialmedia/chat/app/user/usecase"
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	// Follow related use cases and handlers
	followRequestUseCase := followUseCase.NewFollowRequestUseCase(repo)
	followRequestHandler := follow.NewFollowRequestHandler(followRequestUseCase)
	unfollowRequestUseCase := followUseCase.NewUnFollowRequestUseCase(repo)
	unfollowRequestHandler := follow.NewUnFollowRequestHandler(unfollowRequestUseCase)

	blockUserUseCase := followUseCase.NewBlockUserUseCase(repo)
	blockUserHandler := follow.NewBlockUserHandler(blockUserUseCase)
	unblockUserUseCase := followUseCase.NewUnBlockUserUseCase(repo)
	unblockUserHandler := follow.NewUnBlockUserHandler(unblockUserUseCase)

	// User related use cases and handlers
	createUserUseCase := userUseCase.NewCreateUserUseCase(repo)
	updateUserUseCase := userUseCase.NewUpdateUserUseCase(repo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)
	updateUserHandler := user.NewUpdatedUserHandler(updateUserUseCase)

	return map[messaging.MessageType]MessageHandler{
		messaging.UserTypes.UserCreated:     createUserHandler,
		messaging.UserTypes.UserUpdated:     updateUserHandler,
		messaging.UserTypes.UserFollowed:    followRequestHandler,
		messaging.UserTypes.FollowRequest:   followRequestHandler,
		messaging.UserTypes.UnFollowRequest: unfollowRequestHandler,
		messaging.UserTypes.UserBlocked:     blockUserHandler,
		messaging.UserTypes.UserUnBlocked:   unblockUserHandler,
	}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, chatRedisRepo ChatRedisRepository, rabbitMQ Messaging) map[string]interface{} {
	createConversationUseCase := chatUseCase.NewCreateConversationUseCase(repo, rabbitMQ)
	createMessageUseCase := chatUseCase.NewCreateMessageUseCase(repo, chatRedisRepo)
	addParticipantUseCase := chatUseCase.NewAddParticipantUseCase(repo)
	promoteToAdminUseCase := chatUseCase.NewPromoteToAdminUseCase(repo)
	demoteFromAdminUseCase := chatUseCase.NewDemoteFromAdminUseCase(repo)

	createConversationHandler := chat.NewCreateConversationHandler(createConversationUseCase)
	createMessageHandler := chat.NewCreateMessageHandler(createMessageUseCase)
	addParticipantHandler := chat.NewAddParticipantHandler(addParticipantUseCase)
	promoteToAdminHandler := chat.NewPromoteToAdminHandler(promoteToAdminUseCase)
	demoteFromAdminHandler := chat.NewDemoteFromAdminHandler(demoteFromAdminUseCase)

	return map[string]interface{}{

		"createconversation": createConversationHandler,
		"createmessage":      createMessageHandler,
		"addparticipant":     addParticipantHandler,
		"promotetoadmin":     promoteToAdminHandler,
		"demotefromadmin":    demoteFromAdminHandler,
	}
}
func SetupWSHandlers(repo Repository, chatRedisRepo ChatRedisRepository, wsHub Hub) map[string]interface{} {
	chatListenUseCase := chatUseCase.NewChatWebSocketListenUseCase(repo, wsHub)

	chatListenHandler := chat.NewChatWebSocketListenHandler(chatListenUseCase)

	return map[string]interface{}{

		"chatlisten": chatListenHandler,
	}
}
