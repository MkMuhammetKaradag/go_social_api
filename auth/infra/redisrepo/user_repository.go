package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserRedisRepository struct {
	client *redis.Client
}
type UserLogoutNotification struct {
	UserID    uuid.UUID
	Timestamp int64
}

func NewUserRedisRepository(connString, password string, db int) (*UserRedisRepository, error) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &UserRedisRepository{client: RedisClient}, nil
}

func (r *UserRedisRepository) PublishUserLogout(ctx context.Context, userID uuid.UUID) {

	msg := UserLogoutNotification{
		UserID:    userID,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	err = r.client.Publish(ctx, "user:logout", data).Err()
	if err != nil {
		fmt.Println("Redis publish error:", err)
	}

}
func (r *UserRedisRepository) GetRedisClient() *redis.Client {
	return r.client
}
