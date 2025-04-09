package usecase

import (
	"context"
	"time"
)

type RedisRepository interface {
	SetSession(ctx context.Context, key string, userId string, userData map[string]string, expiration time.Duration) error
	GetSession(ctx context.Context, key string) (map[string]string, error)
	DeleteSession(ctx context.Context, key string, userID string) error
}
