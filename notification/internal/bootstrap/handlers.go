package bootstrap

import (
	user "socialmedia/notification/app/user/handler"
	userUseCase "socialmedia/notification/app/user/usecase"
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, repoMongo RepositoryMongo, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	createUserUseCase := userUseCase.NewCreateUserUseCase(repo, repoMongo)
	createUserHandler := user.NewCreatedUserHandler(createUserUseCase)

	return map[messaging.MessageType]MessageHandler{messaging.UserTypes.UserCreated: createUserHandler}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {

	return map[string]interface{}{}
}
func SetupWSHandlers(repo Repository, userRedisRepo UserRedisRepository) map[string]interface{} {

	return map[string]interface{}{}
}
