package initializer

import (
	"context"
	websocket "socialmedia/user/app/user/ws"

	"github.com/redis/go-redis/v9"
)

func InitWebsocket(ctx context.Context, redisRepo UserRedisRepository, client *redis.Client) *websocket.Hub {

	hub := websocket.NewHub(redisRepo, client)
	go hub.Run(ctx)
	return hub
}
