package initializer

import (
	"context"
	"socialmedia/chat/app/chat/ws"

	"github.com/redis/go-redis/v9"
)

func InitWebsocket(ctx context.Context, redisClient *redis.Client, repo ws.Repository) *ws.Hub {

	hub := ws.NewHub(redisClient, repo)
	go hub.Run(ctx)
	return hub
}
