package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(connString, password string, db int) (*RedisRepository, error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &RedisRepository{Client: RedisClient}, nil
}

func (r *RedisRepository) GetSession(ctx context.Context, key string) (map[string]string, error) {
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var userData map[string]string
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return nil, err
	}
	return userData, nil
}
