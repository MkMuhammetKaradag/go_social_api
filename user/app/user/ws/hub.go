package websocket

import (
	"context"
	"fmt"
	"log"
	"socialmedia/user/domain"

	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients     map[uuid.UUID]map[*domain.Client]bool
	register    chan *domain.Client
	unregister  chan *domain.Client
	sessionHub  *SessionHub
	redisReppo  UserRedisRepository
	redisClient *redis.Client
	mutex       sync.RWMutex
}

func NewHub(redisReppo UserRedisRepository, redisClient *redis.Client) *Hub {
	hub := &Hub{
		clients:     make(map[uuid.UUID]map[*domain.Client]bool),
		register:    make(chan *domain.Client),
		unregister:  make(chan *domain.Client),
		redisReppo:  redisReppo,
		redisClient: redisClient,
	}
	hub.sessionHub = NewSessionHub(redisClient, hub)
	return hub
}
func (h *Hub) RegisterClient(client *domain.Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *domain.Client) {
	h.unregister <- client
}

func (h *Hub) Run(ctx context.Context) {
	go func() {
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
	}()
	h.sessionHub.Run(ctx)

}

func (h *Hub) publishStatus(userID uuid.UUID, status string) {
	fmt.Println("user status pıublish userId:", userID, "-status:", status)
	ctx := context.Background()
	h.redisReppo.PublishUserStatus(ctx, userID, status)
}
func (h *Hub) BroadcastToUser(ctx context.Context, userID uuid.UUID, message interface{}) {
	h.mutex.RLock()
	clients, ok := h.clients[userID]
	if !ok {
		h.mutex.RUnlock()
		return
	}

	var clientsToRemove []*domain.Client
	for client := range clients {
		var err error

		switch msg := message.(type) {
		case UserLogoutNotification:

			client.WriteLock.Lock()
			err = client.Conn.WriteJSON(msg)
			client.WriteLock.Unlock()

		default:
			log.Println("Unknown message type")
			continue
		}

		if err != nil {
			log.Println("WebSocket write error:", err)
			client.Conn.Close()
			clientsToRemove = append(clientsToRemove, client)
		}
	}
	h.mutex.RUnlock()

	if len(clientsToRemove) > 0 {
		h.mutex.Lock()
		for _, client := range clientsToRemove {
			delete(h.clients[userID], client)
		}
		h.mutex.Unlock()
	}
}
