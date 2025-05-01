package initializer

import (
	"socialmedia/user/infra/websocket"
)

func InitWebsocket(redisRepo UserRedisRepository) *websocket.Hub {

	hub := websocket.NewHub(redisRepo)
	go hub.Run()
	return hub
}
