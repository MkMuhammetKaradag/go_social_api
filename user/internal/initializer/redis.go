package initializer

import (
	"log"
	"socialmedia/user/infra/redisrepo"
	"socialmedia/user/pkg/config"
)

func InitRedis(appConfig *config.Config) *redisrepo.RedisRepository {
	redisRepo, err := redisrepo.NewRedisRepository(appConfig.Redis.RedisURL, appConfig.Redis.Password, appConfig.Redis.DB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	return redisRepo
}
