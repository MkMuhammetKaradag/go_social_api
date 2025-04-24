package initializer

import (
	"socialmedia/chat/infra/websocket"

	"github.com/redis/go-redis/v9"
)

func InitWebsocket(redisClient *redis.Client) *websocket.Hub {

	hub := websocket.NewHub(redisClient)
	go hub.Run()
	return hub
}
