package websocket

import (
	"context"
	"fmt"
	"socialmedia/user/domain"

	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]map[*domain.Client]bool
	register   chan *domain.Client
	unregister chan *domain.Client
	redisReppo UserRedisRepository
	mutex      sync.RWMutex
}

func NewHub(redisReppo UserRedisRepository) *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*domain.Client]bool),
		register:   make(chan *domain.Client),
		unregister: make(chan *domain.Client),
		redisReppo: redisReppo,
	}
}
func (h *Hub) RegisterClient(client *domain.Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *domain.Client) {
	h.unregister <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if _, ok := h.clients[client.UserID]; !ok {
				h.clients[client.UserID] = make(map[*domain.Client]bool)
			}
			firstConnection := len(h.clients[client.UserID]) == 0
			h.clients[client.UserID][client] = true
			h.mutex.Unlock()

			if firstConnection {
				h.publishStatus(client.UserID, "online")
			}

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients[client.UserID], client)
				if len(h.clients[client.UserID]) == 0 {
					delete(h.clients, client.UserID)
					h.publishStatus(client.UserID, "offline")
				}
				client.Conn.Close()
			}
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) publishStatus(userID uuid.UUID, status string) {
	fmt.Println("user status pıublish userId:", userID, "-status:", status)
	ctx := context.Background()
	h.redisReppo.PublishUserStatus(ctx, userID, status)
}
