package ws

import (
	"context"
	"log"

	"socialmedia/chat/app/chat/usecase"
	"socialmedia/chat/domain"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
)

type RemoveUserFromConversationNotification struct {
	UserID         uuid.UUID `json:"user_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Reason         string    `json:"reason,omitempty"`
	Type           string    `json:"type"` // örnek: "user_removed"
}

// Hub, WebSocket bağlantılarını yöneten ana bileşen
type Hub struct {
	// Tüm aktif istemciler - conversationID -> client -> bool
	clients map[uuid.UUID]map[*domain.Client]bool

	// Sohbet konuşmalarında hangi kullanıcıların olduğunu takip eder
	// conversationID -> userID -> bool
	conversationUsers map[uuid.UUID]map[string]bool

	// İstemci kayıt/silme kanalları
	register   chan *domain.Client
	unregister chan *domain.Client

	// Redis bağlantısı
	redisClient *redis.Client

	// Alt bileşenler
	messageHub  *MessageHub
	statusHub   *StatusHub
	userKickHub *UserKickHub

	// Eşzamanlılık koruması
	mutex sync.RWMutex

	//repository

	repo Repository
}

// NewHub creates a new Hub instance
func NewHub(redisClient *redis.Client, repo Repository) *Hub {
	hub := &Hub{
		clients:           make(map[uuid.UUID]map[*domain.Client]bool),
		conversationUsers: make(map[uuid.UUID]map[string]bool),
		register:          make(chan *domain.Client),
		unregister:        make(chan *domain.Client),
		redisClient:       redisClient,
		repo:              repo,
	}

	// Alt bileşenleri oluştur
	hub.messageHub = NewMessageHub(redisClient, hub)
	hub.statusHub = NewStatusHub(redisClient, hub)
	hub.userKickHub = NewUserKickHub(redisClient, hub)

	return hub
}

// Run başlatır tüm Hub aktivitelerini
func (h *Hub) Run(ctx context.Context) {
	// Ana hub döngüsü
	go func() {
		for {
			select {
			case client := <-h.register:
				h.registerClient(client)
			case client := <-h.unregister:
				h.unregisterClient(client)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Alt bileşenleri başlat
	go h.messageHub.Run(ctx)
	go h.statusHub.Run(ctx)
	go h.userKickHub.Run(ctx)
}

// RegisterClient bir WebSocket bağlantısını Hub'a kaydeder ve kullanıcıyı sohbetle ilişkilendirir
func (h *Hub) RegisterClient(client *domain.Client, userID uuid.UUID) {
	h.register <- client

	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Kullanıcıyı sohbetle ilişkilendir
	if _, ok := h.conversationUsers[client.ConversationID]; !ok {
		h.conversationUsers[client.ConversationID] = make(map[string]bool)
	}
	h.conversationUsers[client.ConversationID][userID.String()] = true

	// Kullanıcıya mevcut durumları gönder
	go h.statusHub.SendInitialUserStatuses(client, client.ConversationID)
}

// UnregisterClient removes a client from the hub
func (h *Hub) UnregisterClient(client *domain.Client, userID uuid.UUID) {
	h.unregister <- client

	// İlişkilendirmeyi kaldırma işlemini unregisterClient içinde yapacağız
}

// registerClient handles client registration (internal)
func (h *Hub) registerClient(client *domain.Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client.ConversationID]; !ok {
		h.clients[client.ConversationID] = make(map[*domain.Client]bool)
	}
	h.clients[client.ConversationID][client] = true
}

// unregisterClient handles client unregistration (internal)
func (h *Hub) unregisterClient(client *domain.Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client.ConversationID]; ok {
		delete(h.clients[client.ConversationID], client)
		if len(h.clients[client.ConversationID]) == 0 {
			delete(h.clients, client.ConversationID)
		}

		// Bağlantıyı kapat
		client.Conn.Close()
	}
}

// LoadConversationMembers sohbetin katılımcılarını yükler
func (h *Hub) LoadConversationMembers(ctx context.Context, conversationID uuid.UUID, repo usecase.Repository) error {
	// Sohbet katılımcılarını veritabanından çek
	participants, err := h.repo.GetParticipants(ctx, conversationID)
	if err != nil {
		return err
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.conversationUsers[conversationID]; !ok {
		h.conversationUsers[conversationID] = make(map[string]bool)
	}

	// Tüm katılımcıları kaydet
	for _, participant := range participants {
		h.conversationUsers[conversationID][participant.String()] = true
	}

	return nil
}

// IsConversationLoaded sohbetin yüklenip yüklenmediğini kontrol eder
func (h *Hub) IsConversationLoaded(conversationID uuid.UUID) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, exists := h.conversationUsers[conversationID]
	return exists
}

// GetConversationUsers, bir sohbetteki kullanıcı kimliklerini döndürür
func (h *Hub) GetConversationUsers(conversationID uuid.UUID) map[string]bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if users, ok := h.conversationUsers[conversationID]; ok {
		// Kopya oluştur, doğrudan referans değil
		result := make(map[string]bool, len(users))
		for userID, val := range users {
			result[userID] = val
		}
		return result
	}

	return make(map[string]bool)
}

func (h *Hub) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) bool {
	exists, err := h.repo.IsBlocked(ctx, blockerID, blockedID)
	if err != nil {
		log.Println("repo block check error:", err)
		return false
	}
	return exists
}

func (h *Hub) KickUserFromConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) {
	h.mutex.Lock()
	clients, ok := h.clients[conversationID]
	if !ok {
		h.mutex.Unlock()
		return
	}

	var targetClient *domain.Client
	for client := range clients {
		if client.UserID == userID {
			targetClient = client
			break
		}
	}
	if targetClient == nil {
		h.mutex.Unlock()
		return
	}
	h.mutex.Unlock()

	// Kullanıcıya önce bildirim gönder
	h.RemoveUserFromConversation(ctx, conversationID, userID)

	// WebSocket bağlantısını düzgün kapat
	targetClient.WriteLock.Lock()
	_ = targetClient.Conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "You have been removed from the conversation by admin."),
	)
	targetClient.WriteLock.Unlock()
	targetClient.Conn.Close()

	// Haritadan client'ı sil
	h.mutex.Lock()
	delete(h.clients[conversationID], targetClient)
	if users, exists := h.conversationUsers[conversationID]; exists {
		delete(users, userID.String())
	}
	h.mutex.Unlock()
}

// RemoveUserFromConversation, belirli bir kullanıcıyı belirtilen sohbetten çıkarır

func (h *Hub) RemoveUserFromConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) {
	var shouldBroadcast bool
	msg := RemoveUserFromConversationNotification{
		ConversationID: conversationID,
		UserID:         userID,
		Type:           "user_removed",
	}

	h.mutex.Lock()
	if users, exists := h.conversationUsers[conversationID]; exists {
		if _, ok := users[userID.String()]; ok {
			delete(users, userID.String())
			log.Printf("User %s removed from conversation %s\n", userID, conversationID)
			shouldBroadcast = true

			if len(users) == 0 {
				delete(h.conversationUsers, conversationID)
			}
		}
	}
	h.mutex.Unlock()

	if shouldBroadcast {
		h.BroadcastToConversation(ctx, conversationID, msg)
	}
}

func (h *Hub) BroadcastToConversation(ctx context.Context, conversationID uuid.UUID, message interface{}) {
	h.mutex.RLock()
	clients, ok := h.clients[conversationID]
	if !ok {
		h.mutex.RUnlock()
		return
	}

	var clientsToRemove []*domain.Client
	for client := range clients {
		var err error

		switch msg := message.(type) {
		case MessageNotification:
			if h.IsBlocked(ctx, client.UserID, msg.UserID) {
				continue
			}
			client.WriteLock.Lock()
			err = client.Conn.WriteJSON(msg)
			client.WriteLock.Unlock()

		case UserStatusNotification, RemoveUserFromConversationNotification:
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

	// Silme işlemleri ayrı Lock ile yapılmalı
	if len(clientsToRemove) > 0 {
		h.mutex.Lock()
		for _, client := range clientsToRemove {
			delete(h.clients[conversationID], client)
		}
		h.mutex.Unlock()
	}
}
