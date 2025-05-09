// status_hub.go
package ws

import (
	"context"
	"encoding/json"
	"log"
	"socialmedia/chat/domain"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// UserStatusNotification, kullanıcı durum bildirimlerinin yapısı
type UserStatusNotification struct {
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"` // "online" veya "offline"
	Timestamp int64     `json:"timestamp"`
	Type      string    `json:"type"` // Bildirim türü, "user_status" olarak sabit
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

func (u UserStatusNotification) GetUserID() uuid.UUID {
	return u.UserID
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
			convID := conversationID
			// if err != nil {
			// 	log.Println("Invalid conversation ID:", err)
			// 	continue
			// }

			// Sohbetteki tüm istemcilere yayınla
			sh.parentHub.BroadcastToConversation(ctx, convID, statusNotif)
		}
	}
}

// getConversationsWithUser, bir kullanıcının bulunduğu tüm sohbetleri döndürür
func (sh *StatusHub) getConversationsWithUser(userID uuid.UUID) map[uuid.UUID]bool {
	result := make(map[uuid.UUID]bool)

	// Ana Hub'dan tüm sohbetleri ve kullanıcılarını al
	for conversationID, users := range sh.getAllConversationUsers() {
		if _, ok := users[userID]; ok {
			result[conversationID] = true
		}
	}

	return result
}

// getAllConversationUsers, tüm sohbetlerdeki kullanıcıları döndürür (yardımcı metod)
func (sh *StatusHub) getAllConversationUsers() map[uuid.UUID]map[uuid.UUID]UserInfo {
	result := make(map[uuid.UUID]map[uuid.UUID]UserInfo)

	sh.parentHub.mutex.RLock()
	defer sh.parentHub.mutex.RUnlock()

	for convID, users := range sh.parentHub.conversationUsers {
		usersCopy := make(map[uuid.UUID]UserInfo, len(users))
		for userID, info := range users {
			usersCopy[userID] = info
		}
		result[convID] = usersCopy
	}

	return result
}

// SendInitialUserStatuses, yeni bağlanan istemciye tüm kullanıcı durumlarını gönderir
func (sh *StatusHub) SendInitialUserStatuses(client *domain.Client, conversationID uuid.UUID) {
	users := sh.parentHub.GetConversationUsers(conversationID)
	if len(users) == 0 {
		return
	}

	statusUpdates := []UserStatusNotification{}

	for userID, info := range users {
		if userID == client.UserID {
			continue
		}

		status := UserStatusNotification{
			UserID:    userID,
			Username:  info.Username,
			Avatar:    info.Avatar,
			Timestamp: time.Now().Unix(),
			Type:      "user_status",
		}

		// Kullanıcının durumu
		statusVal, exists := sh.userStatuses.Load(userID)
		if exists {
			status.Status = statusVal.(UserStatusNotification).Status
		} else {
			// Redis'ten çek
			statusStr, err := sh.redisClient.Get(context.Background(), "user:status:"+userID.String()).Result()
			if err == redis.Nil {
				status.Status = "offline"
			} else if err != nil {
				log.Println("Redis get error:", err)
				status.Status = "unknown"
			} else {
				var redisStatus UserStatusNotification
				err = json.Unmarshal([]byte(statusStr), &redisStatus)
				if err != nil {
					log.Println("Redis status unmarshal error:", err)
					status.Status = "unknown"
				} else {
					status.Status = redisStatus.Status
					sh.userStatuses.Store(userID, status)
				}
			}
		}

		statusUpdates = append(statusUpdates, status)
	}

	// Burada `statusUpdates` client'a gönderilebilir
	// Örnek:
	payload, _ := json.Marshal(statusUpdates)
	client.WriteLock.Lock()
	_ = client.Conn.WriteMessage(websocket.TextMessage, payload)
	client.WriteLock.Unlock()
}

// UpdateUserStatus, bir kullanıcının durumunu günceller ve yayınlar
func (sh *StatusHub) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error {
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
	err = sh.redisClient.Set(ctx, "user:status:"+userID.String(), string(data), 24*time.Hour).Err()
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
