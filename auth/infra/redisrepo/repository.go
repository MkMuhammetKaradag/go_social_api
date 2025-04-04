package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
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

func (r *RedisRepository) SetSession(ctx context.Context, key string, userData map[string]string, expiration time.Duration) error {
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, jsonData, expiration).Err()
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

func (r *RedisRepository) DeleteSession(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// func (r *RedisRepository) PublishStatus(ctx context.Context, userID string, status string) error {
// 	return r.Client.Publish(ctx, "user_status", userID+":"+status).Err()
// }

// func (r *RedisRepository) PublishChatMessage(ctx context.Context, chatID string, content string, senderID string) error {
// 	return r.Client.Publish(ctx, "send_Message", chatID+":"+content+":"+senderID).Err()
// }
