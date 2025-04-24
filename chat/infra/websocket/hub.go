package websocket

import (
	"context"
	"encoding/json"
	"log"
	"socialmedia/chat/domain"
	"sync"

	"github.com/google/uuid"

	// "github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type MessageNotification struct {
	MessageID      uuid.UUID        `json:"message_id"`
	ConversationID uuid.UUID        `json:"conversation_id"`
	SenderID       uuid.UUID        `json:"sender_id"`
	Content        string           `json:"content"`
	CreatedAt      string           `json:"created_at"`
	HasAttachments bool             `json:"has_attachments"`
	Attachments    []AttachmentInfo `json:"attachments,omitempty"`
}

type AttachmentInfo struct {
	ID       uuid.UUID `json:"id"`
	FileURL  string    `json:"file_url"`
	FileType string    `json:"file_type"`
}


type Hub struct {
	clients     map[uuid.UUID]map[*domain.Client]bool
	register    chan *domain.Client
	redisClient *redis.Client
	unregister  chan *domain.Client
	mutex       sync.RWMutex
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:     make(map[uuid.UUID]map[*domain.Client]bool),
		register:    make(chan *domain.Client),
		unregister:  make(chan *domain.Client),
		redisClient: redisClient,
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
			if _, ok := h.clients[client.ConversationID]; !ok {
				h.clients[client.ConversationID] = make(map[*domain.Client]bool)
			}
			h.clients[client.ConversationID][client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.ConversationID]; ok {
				delete(h.clients[client.ConversationID], client)
				if len(h.clients[client.ConversationID]) == 0 {
					delete(h.clients, client.ConversationID)
				}
				client.Conn.Close()
			}
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) ListenRedisSendMessage(ctx context.Context, channelName string) {
	pubsub := h.redisClient.Subscribe(ctx, channelName)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("Redis sub error:", err)
			continue
		}

		var notification MessageNotification
		err = json.Unmarshal([]byte(msg.Payload), &notification)
		if err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		conversationID := notification.ConversationID

		h.mutex.RLock()
		if clients, ok := h.clients[conversationID]; ok {
			for client := range clients {
				client.WriteLock.Lock()
				err := client.Conn.WriteJSON(notification)
				client.WriteLock.Unlock()

				if err != nil {
					log.Println("WebSocket write error:", err)
					client.Conn.Close()
					delete(h.clients[conversationID], client)
				}
			}
		}
		h.mutex.RUnlock()
	}
}
