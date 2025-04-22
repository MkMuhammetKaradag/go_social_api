package bootstrap

import (
	"socialmedia/shared/messaging"
	follow "socialmedia/user/app/follow/handler"
	followUseCase "socialmedia/user/app/follow/usecase"
	user "socialmedia/user/app/user/handler"
	userUseCase "socialmedia/user/app/user/usecase"
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
	profileUseCase := userUseCase.NewProfileUseCase(redisRepo, repo)
	updateUserUseCase := userUseCase.NewUpdateUserUseCase(redisRepo, repo)
	getUserUseCase := userUseCase.NewGetUserUseCase(redisRepo, repo)
	searchUsersUseCase := userUseCase.NewSearchUserUseCase(redisRepo, repo)
	updateAvatarUseCase := userUseCase.NewUpdateAvatarUseCase(redisRepo, repo)
	updateBanerUseCase := userUseCase.NewUpdateBannerUseCase(redisRepo, repo)

	profileUserHandler := user.NewProfileUserHandler(profileUseCase)
	updateUserHandler := user.NewUpdateUserHandler(updateUserUseCase)
	getUserHandler := user.NewGetUserHandler(getUserUseCase)
	searchUsersHandler := user.NewSearchUserHandler(searchUsersUseCase)
	updateAvatarHandler := user.NewUpdateAvatarHandler(updateAvatarUseCase)
	updateBannerHandler := user.NewUpdateBannerHandler(updateBanerUseCase)

	return map[string]interface{}{
		"profile":     profileUserHandler,
		"update":      updateUserHandler,
		"getUser":     getUserHandler,
		"searchusers": searchUsersHandler,
		"avatar":      updateAvatarHandler,
		"banner":      updateBannerHandler,
	}
}
