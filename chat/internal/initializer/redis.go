package initializer

import (
	"log"
	"socialmedia/chat/infra/redisrepo"
	"socialmedia/chat/pkg/config"
)

func InitRedis(appConfig config.Config) *redisrepo.RedisRepository {
	redisRepo, err := redisrepo.NewRedisRepository(appConfig.Redis.RedisURL, appConfig.Redis.Password, appConfig.Redis.SessionDB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	return redisRepo
}
func InitChatRedis(appConfig config.Config) *redisrepo.ChatRedisRepository {
	redisRepo, err := redisrepo.NewChatRedisRepository(appConfig.Redis.RedisURL, appConfig.Redis.Password, appConfig.Redis.ChatDB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	return redisRepo
}
