package initializer

import (
	websocket "socialmedia/user/app/user/ws"
)

func InitWebsocket(redisRepo UserRedisRepository) *websocket.Hub {

	hub := websocket.NewHub(redisRepo)
	go hub.Run()
	return hub
}
