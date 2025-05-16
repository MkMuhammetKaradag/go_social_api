package initializer

import (
	"log"
	"socialmedia/auth/infra/redisrepo"
	"socialmedia/auth/pkg/config"
)

func InitRedis(appConfig config.Config) *redisrepo.RedisRepository {
	redisRepo, err := redisrepo.NewRedisRepository(appConfig.Redis.RedisURL, appConfig.Redis.Password, appConfig.Redis.DB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	return redisRepo
}
func InitUserRedis(appConfig config.Config) *redisrepo.UserRedisRepository {
	redisRepo, err := redisrepo.NewUserRedisRepository(appConfig.Redis.RedisURL, appConfig.Redis.Password, appConfig.Redis.UserDB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	return redisRepo
}
