// status_hub.go
package ws

import (
	"context"
	"encoding/json"
	"log"
	"socialmedia/chat/domain"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// UserStatusNotification, kullanıcı durum bildirimlerinin yapısı
type UserStatusNotification struct {
	UserID    string `json:"user_id"`
	Status    string `json:"status"` // "online" veya "offline"
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"` // Bildirim türü, "user_status" olarak sabit
}

// StatusHub, kullanıcı durumu işlemlerini yöneten bileşen
type StatusHub struct {
	redisClient *redis.Client
	parentHub   *Hub

	// Kullanıcı durumları için önbellek
	userStatuses sync.Map

	// Periyodik temizleme için zamanlayıcı
	cleanupTicker *time.Ticker
}

// NewStatusHub, yeni bir StatusHub örneği oluşturur
func NewStatusHub(redisClient *redis.Client, parent *Hub) *StatusHub {
	return &StatusHub{
		redisClient:   redisClient,
		parentHub:     parent,
		userStatuses:  sync.Map{},
		cleanupTicker: time.NewTicker(30 * time.Minute), // 30 dakikada bir eskimiş durumları temizle
	}
}

// Run, durum dinleme ve işleme işlevini başlatır
func (sh *StatusHub) Run(ctx context.Context) {
	// Redis'ten kullanıcı durumlarını dinle
	go sh.listenForStatusUpdates(ctx)

	// Önbellek temizleme işlemini başlat
	go sh.cleanupCacheRoutine(ctx)
}

// listenForStatusUpdates, Redis'ten gelen kullanıcı durumu güncellemelerini dinler
func (sh *StatusHub) listenForStatusUpdates(ctx context.Context) {
	pubsub := sh.redisClient.Subscribe(ctx, "user:status")
	defer pubsub.Close()

	// Kanal dinleme döngüsü
	for {
		select {
		case <-ctx.Done():
			// Bağlam iptal edildiğinde çık
			return

		default:
			// Redis'ten mesaj al
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Println("Redis status subscription error:", err)
				continue
			}

			// Durumu işle
			var statusNotif UserStatusNotification
			err = json.Unmarshal([]byte(msg.Payload), &statusNotif)
			if err != nil {
				log.Println("Status unmarshal error:", err)
				continue
			}

			// Bildirim türünü belirt
			statusNotif.Type = "user_status"

			// Önbelleği güncelle
			sh.userStatuses.Store(statusNotif.UserID, statusNotif)

			// İlgili sohbetlere durumu bildir
			sh.broadcastStatusToRelevantConversations(ctx, statusNotif)
		}
	}
}

// broadcastStatusToRelevantConversations, durum güncellemesini ilgili sohbetlere yayınlar
func (sh *StatusHub) broadcastStatusToRelevantConversations(ctx context.Context, statusNotif UserStatusNotification) {
	// Kullanıcının bulunduğu tüm sohbetleri bul
	// parentHub üzerinden tüm sohbetleri tara ve kullanıcının olduğu sohbetleri bul
	for conversationID, users := range sh.getConversationsWithUser(statusNotif.UserID) {
		if users {
			convID, err := uuid.Parse(conversationID)
			if err != nil {
				log.Println("Invalid conversation ID:", err)
				continue
			}

			// Sohbetteki tüm istemcilere yayınla
			sh.parentHub.BroadcastToConversation(ctx, convID, statusNotif)
		}
	}
}

// getConversationsWithUser, bir kullanıcının bulunduğu tüm sohbetleri döndürür
func (sh *StatusHub) getConversationsWithUser(userID string) map[string]bool {
	result := make(map[string]bool)

	// Ana Hub'dan tüm sohbetleri ve kullanıcılarını al
	for conversationID, users := range sh.getAllConversationUsers() {
		if _, ok := users[userID]; ok {
			result[conversationID.String()] = true
		}
	}

	return result
}

// getAllConversationUsers, tüm sohbetlerdeki kullanıcıları döndürür (yardımcı metod)
func (sh *StatusHub) getAllConversationUsers() map[uuid.UUID]map[string]bool {
	result := make(map[uuid.UUID]map[string]bool)

	// Burada Ana Hub'dan conversationUsers map'ini kopyalamak gerekiyor
	// Ancak bu örnek için basitleştirmek adına doğrudan erişim kullanılıyor
	sh.parentHub.mutex.RLock()
	defer sh.parentHub.mutex.RUnlock()

	for convID, users := range sh.parentHub.conversationUsers {
		usersCopy := make(map[string]bool, len(users))
		for userID, val := range users {
			usersCopy[userID] = val
		}
		result[convID] = usersCopy
	}

	return result
}

// SendInitialUserStatuses, yeni bağlanan istemciye tüm kullanıcı durumlarını gönderir
func (sh *StatusHub) SendInitialUserStatuses(client *domain.Client, conversationID uuid.UUID) {
	// Sohbetteki tüm kullanıcıları al
	users := sh.parentHub.GetConversationUsers(conversationID)
	if len(users) == 0 {
		return
	}

	// Tüm kullanıcıların durumlarını topla
	statusUpdates := []UserStatusNotification{}

	for userID := range users {
		// Kullanıcının durumunu önbellekten al
		statusVal, exists := sh.userStatuses.Load(userID)
		var status UserStatusNotification

		if !exists {
			// Önbellekte yoksa Redis'ten çek
			statusStr, err := sh.redisClient.Get(context.Background(), "user:status:"+userID).Result()
			if err == redis.Nil {
				// Varsayılan olarak offline
				status = UserStatusNotification{
					UserID:    userID,
					Status:    "offline",
					Timestamp: time.Now().Unix(),
					Type:      "user_status",
				}
			} else if err != nil {
				log.Println("Redis get error:", err)
				status = UserStatusNotification{
					UserID:    userID,
					Status:    "unknown",
					Timestamp: time.Now().Unix(),
					Type:      "user_status",
				}
			} else {
				// Redis'ten gelen string'i ayrıştır
				var redisStatus UserStatusNotification
				err = json.Unmarshal([]byte(statusStr), &redisStatus)
				if err != nil {
					log.Println("Redis status unmarshal error:", err)
					status = UserStatusNotification{
						UserID:    userID,
						Status:    "unknown",
						Timestamp: time.Now().Unix(),
						Type:      "user_status",
					}
				} else {
					status = redisStatus
					status.Type = "user_status"
				}

				// Önbelleğe kaydet
				sh.userStatuses.Store(userID, status)
			}
		} else {
			// Önbellekten al
			status = statusVal.(UserStatusNotification)
		}

		statusUpdates = append(statusUpdates, status)
	}

	// Tüm kullanıcı durumlarını tek bir mesajda gönder
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

// UpdateUserStatus, bir kullanıcının durumunu günceller ve yayınlar
func (sh *StatusHub) UpdateUserStatus(ctx context.Context, userID string, status string) error {
	// Durum bildirimini oluştur
	statusNotif := UserStatusNotification{
		UserID:    userID,
		Status:    status,
		Timestamp: time.Now().Unix(),
		Type:      "user_status",
	}

	// Önbelleği güncelle
	sh.userStatuses.Store(userID, statusNotif)

	// Redis'e yayınla
	data, err := json.Marshal(statusNotif)
	if err != nil {
		return err
	}

	// Ayrıca Redis'e kalıcı olarak kaydet
	err = sh.redisClient.Set(ctx, "user:status:"+userID, string(data), 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return sh.redisClient.Publish(ctx, "user:status", string(data)).Err()
}

// cleanupCacheRoutine, önbellekteki eski durumları periyodik olarak temizler
func (sh *StatusHub) cleanupCacheRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			sh.cleanupTicker.Stop()
			return
		case <-sh.cleanupTicker.C:
			sh.cleanupCache()
		}
	}
}

// cleanupCache, eski kullanıcı durumlarını önbellekten temizler (24 saatten eski)
func (sh *StatusHub) cleanupCache() {
	now := time.Now().Unix()
	thresholdTime := now - 24*60*60 // 24 saat

	// Tüm durumları kontrol et ve eski olanları temizle
	sh.userStatuses.Range(func(key, value interface{}) bool {
		if status, ok := value.(UserStatusNotification); ok {
			if status.Timestamp < thresholdTime {
				sh.userStatuses.Delete(key)
			}
		}
		return true
	})
}
