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
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)

	return map[messaging.MessageType]MessageHandler{
		messaging.UserTypes.UserCreated:     createUserHandler,
		messaging.UserTypes.UserFollowed:    followRequestHandler,
		messaging.UserTypes.FollowRequest:   followRequestHandler,
		messaging.UserTypes.UnFollowRequest: unfollowRequestHandler,
		messaging.UserTypes.UserBlocked:     blockUserHandler,
		messaging.UserTypes.UserUnBlocked:   unblockUserHandler,
	}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {
	createConversationUseCase := chatUseCase.NewCreateConversationUseCase(repo)
	createConversationHandler := chat.NewCreateConversationHandler(createConversationUseCase)

	return map[string]interface{}{

		"createconversation": createConversationHandler,
	}
}
