package websocket

import (
	"context"
	"encoding/json"
	"log"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/chat/domain"
	"sync"

	"github.com/google/uuid"

	// "github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type UserStatusNotification struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"` // "online" veya "offline"
	Timestamp int64  `json:"timestamp"`
}

var userStatusCache = sync.Map{}

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
	// Yeni eklenen alan:
	conversationUsers map[uuid.UUID]map[string]bool // conversationID -> map[userID]bool
}

// Hub constructor'ını güncelleyin
func NewHub(redisClient *redis.Client) *Hub {
	hub := &Hub{
		clients:           make(map[uuid.UUID]map[*domain.Client]bool),
		register:          make(chan *domain.Client),
		unregister:        make(chan *domain.Client),
		redisClient:       redisClient,
		conversationUsers: make(map[uuid.UUID]map[string]bool),
	}

	// Kullanıcı durumlarını dinlemeye başla
	go hub.listenUserStatusUpdates()

	return hub
}

// RegisterClient metodunu güncelleyin - kullanıcıları sohbetle ilişkilendirin
func (h *Hub) RegisterClient(client *domain.Client, userID uuid.UUID) {
	// Mevcut register işlemi
	h.register <- client

	// Kullanıcıyı sohbetle ilişkilendir
	h.mutex.Lock()
	if _, ok := h.conversationUsers[client.ConversationID]; !ok {
		h.conversationUsers[client.ConversationID] = make(map[string]bool)
	}
	h.conversationUsers[client.ConversationID][userID.String()] = true
	h.mutex.Unlock()
}

// UnregisterClient metodunu güncelleyin
func (h *Hub) UnregisterClient(client *domain.Client, userID uuid.UUID) {
	// Mevcut unregister işlemi
	h.unregister <- client

	// İlişkiyi kaldır (eğer başka bağlantısı yoksa)
	h.mutex.Lock()
	if _, ok := h.conversationUsers[client.ConversationID]; ok {
		// Kullanıcının aynı sohbette başka bağlantısı var mı kontrol et
		hasOtherConnections := false
		if clients, exists := h.clients[client.ConversationID]; exists {
			for otherClient := range clients {
				if otherClient != client {
					// Burada client userID'sini kontrol etmeniz gerekebilir
					// Bu örnekte basit tutmak için atlandı
					hasOtherConnections = true
					break
				}
			}
		}

		if !hasOtherConnections {
			delete(h.conversationUsers[client.ConversationID], userID.String())
			if len(h.conversationUsers[client.ConversationID]) == 0 {
				delete(h.conversationUsers, client.ConversationID)
			}
		}
	}
	h.mutex.Unlock()
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

func (h *Hub) listenUserStatusUpdates() {
	pubsub := h.redisClient.Subscribe(context.Background(), "user:status")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Println("Redis user status sub error:", err)
			continue
		}

		var statusNotif UserStatusNotification
		err = json.Unmarshal([]byte(msg.Payload), &statusNotif)
		if err != nil {
			log.Println("User status unmarshal error:", err)
			continue
		}

		// Önbelleği güncelle
		userStatusCache.Store(statusNotif.UserID, statusNotif.Status)

		// İlgili kullanıcının bulunduğu tüm sohbetlere durumu bildir
		h.broadcastUserStatusToRelevantConversations(statusNotif)
	}
}

// Kullanıcı durumunu ilgili sohbetlere yayınlayan yeni metod
func (h *Hub) broadcastUserStatusToRelevantConversations(statusNotif UserStatusNotification) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Kullanıcının hangi sohbetlerde olduğunu bul
	for conversationID, users := range h.conversationUsers {
		if _, exists := users[statusNotif.UserID]; exists {
			// Bu sohbetteki tüm bağlı istemcilere user status bilgisini gönder
			if clients, ok := h.clients[conversationID]; ok {
				for client := range clients {
					client.WriteLock.Lock()
					err := client.Conn.WriteJSON(map[string]interface{}{
						"type":      "user_status",
						"user_id":   statusNotif.UserID,
						"status":    statusNotif.Status,
						"timestamp": statusNotif.Timestamp,
					})
					client.WriteLock.Unlock()

					if err != nil {
						log.Println("WebSocket status write error:", err)
					}
				}
			}
		}
	}
}
func (h *Hub) LoadConversationMembers(ctx context.Context, conversationID uuid.UUID, repo usecase.Repository) error {
	// Sohbet katılımcılarını veritabanından çek
	participants, err := repo.GetParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	h.mutex.Lock()
	if _, ok := h.conversationUsers[conversationID]; !ok {
		h.conversationUsers[conversationID] = make(map[string]bool)
	}

	// Tüm katılımcıları kaydet
	for _, participant := range participants {
		h.conversationUsers[conversationID][participant.String()] = true
	}
	h.mutex.Unlock()

	return nil
}

// Sohbete bağlanan kullanıcıya tüm kullanıcıların durumlarını gönder
func (h *Hub) SendInitialUserStatuses(client *domain.Client, conversationID uuid.UUID) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if users, ok := h.conversationUsers[conversationID]; ok {
		statusUpdates := []map[string]interface{}{}

		for userID := range users {
			// Kullanıcının durumunu önbellekten veya Redis'ten al
			status, exists := userStatusCache.Load(userID)
			if !exists {
				// Önbellekte yoksa Redis'ten çek
				statusStr, err := h.redisClient.Get(context.Background(), "user:status:"+userID).Result()
				if err == redis.Nil {
					status = "offline" // Varsayılan olarak offline
				} else if err != nil {
					log.Println("Redis get error:", err)
					status = "unknown"
				} else {
					status = statusStr
					// Önbelleğe kaydet
					userStatusCache.Store(userID, status)
				}
			}

			statusUpdates = append(statusUpdates, map[string]interface{}{
				"type":    "user_status",
				"user_id": userID,
				"status":  status,
			})
		}

		// Tüm kullanıcıların durumlarını tek bir mesajda gönder
		if len(statusUpdates) > 0 {
			client.WriteLock.Lock()
			err := client.Conn.WriteJSON(map[string]interface{}{
				"type":           "initial_user_statuses",
				"status_updates": statusUpdates,
			})
			client.WriteLock.Unlock()

			if err != nil {
				log.Println("Error sending initial statuses:", err)
			}
		}
	}
}
