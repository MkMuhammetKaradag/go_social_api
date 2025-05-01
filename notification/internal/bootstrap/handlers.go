package bootstrap

import (
	"socialmedia/shared/messaging"
)

func SetupMessageHandlers(repo Repository, redisRepo RedisRepository) map[messaging.MessageType]MessageHandler {
	// Follow related use cases and handlers

	return map[messaging.MessageType]MessageHandler{}
}

func SetupHTTPHandlers(repo Repository, redisRepo RedisRepository, rabbitMQ Messaging) map[string]interface{} {

	return map[string]interface{}{}
}
func SetupWSHandlers(repo Repository, userRedisRepo UserRedisRepository) map[string]interface{} {

	return map[string]interface{}{}
}
