package bootstrap

import (
	follow "socialmedia/follow/app/follow/handler"
	followUseCase "socialmedia/follow/app/follow/usecase"
	user "socialmedia/follow/app/user/handler"
	userUseCase "socialmedia/follow/app/user/usecase"
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	createUserUseCase := userUseCase.NewCreateUserUseCase(repo)
	updateUserUseCase := userUseCase.NewUpdateUserUseCase(repo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)
	updateUserHandler := user.NewUpdatedUserHandler(updateUserUseCase)
	return map[messaging.MessageType]MessageHandler{
		messaging.UserTypes.UserCreated: createUserHandler,
		messaging.UserTypes.UserUpdated: updateUserHandler,
	}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {
	followRequestUseCase := followUseCase.NewFollowRequestUseCase(redisRepo, repo, rabbitMQ)
	unfollowRequestUseCase := followUseCase.NewUnFollowRequestUseCase(redisRepo, repo, rabbitMQ)
	blockUserUseCase := followUseCase.NewBlockUserUseCase(redisRepo, repo, rabbitMQ)
	unblockUserUseCase := followUseCase.NewUnblockUserUseCase(redisRepo, repo, rabbitMQ)
	incomingRequestUseCase := followUseCase.NewIncomingRequestsUseCase(redisRepo, repo)
	outgoingRequestUseCase := followUseCase.NewOutgoingRequestsUseCase(redisRepo, repo)
	getBlockedUsersUseCase := followUseCase.NewGetBlockedUsersUseCase(redisRepo, repo)
	acceptFollowRequestUseCase := followUseCase.NewAcceptFollowRequestUseCase(redisRepo, repo, rabbitMQ)
	rejectFollowRequestUseCase := followUseCase.NewRejectFollowRequestUseCase(redisRepo, repo, rabbitMQ)

	fallawRequestHandler := follow.NewFollowRequestHandler(followRequestUseCase)
	unfallawRequestHandler := follow.NewUnFollowRequestHandler(unfollowRequestUseCase)
	blockUserHandler := follow.NewBlockUserHandler(blockUserUseCase)
	unblockUserHandler := follow.NewUnblockUserHandler(unblockUserUseCase)
	incomingRequestHandler := follow.NewIncomingRequestsHandler(incomingRequestUseCase)
	outgoingRequestHandler := follow.NewOutgoingRequestsHandler(outgoingRequestUseCase)
	getBlockedUsersHandler := follow.NewGetBlockedUsersHandler(getBlockedUsersUseCase)
	acceptFollowHandler := follow.NewAcceptFollowRequestHandler(acceptFollowRequestUseCase)
	rejectFollowHandler := follow.NewRejectFollowRequestHandler(rejectFollowRequestUseCase)

	return map[string]interface{}{
		"follow":          fallawRequestHandler,
		"unfollow":        unfallawRequestHandler,
		"block":           blockUserHandler,
		"unblock":         unblockUserHandler,
		"incomingRequest": incomingRequestHandler,
		"outgoingRequest": outgoingRequestHandler,
		"getBlockedUsers": getBlockedUsersHandler,
		"accept":          acceptFollowHandler,
		"reject":          rejectFollowHandler,
	}
}
