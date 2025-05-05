package mongodb

import (
	"context"
	"fmt"
	"log"
	"socialmedia/notification/domain"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewRepository(connString, dbName string) (*Repository, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to verify the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")

	database := client.Database(dbName)

	repo := &Repository{
		client:   client,
		database: database,
	}

	return repo, nil
}

// Koleksiyon adını dinamik olarak alabilen fonksiyon
func (r *Repository) GetCollection(collectionName string) *mongo.Collection {
	return r.database.Collection(collectionName)
}

type User struct {
	ID        string    `bson:"id"` // UUID as string
	Username  string    `bson:"username"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

// CreateUser fonksiyonu
func (r *Repository) CreateUser(ctx context.Context, id uuid.UUID, username string) error {
	// Kullanıcı koleksiyonuna erişiyoruz
	usersCollection := r.GetCollection("users")

	// UUID'yi string formatına dönüştür
	userId := id.String()

	// Şu anki zamanı al
	now := time.Now()

	// Veritabanında kullanıcıyı oluşturma ya da güncelleme işlemi
	filter := bson.M{"id": userId}

	// Kullanıcının mevcut olup olmadığını kontrol et
	var existingUser User
	err := usersCollection.FindOne(ctx, filter).Decode(&existingUser)

	update := bson.M{
		"$set": bson.M{
			"id":        userId,
			"username":  username,
			"updatedAt": now,
		},
	}

	// Eğer kullanıcı yoksa createdAt ekle
	if err == mongo.ErrNoDocuments {
		update["$set"].(bson.M)["createdAt"] = now
	}

	options := options.Update().SetUpsert(true) // Upsert: varsa güncelle, yoksa ekle

	_, err = usersCollection.UpdateOne(ctx, filter, update, options)
	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}

	return nil
}
func (r *Repository) CreateNotification(ctx context.Context, notification domain.Notification) error {
	// Bildirim ID yoksa oluştur
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}

	// Zaman bilgilerini ayarla
	now := time.Now()
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = now
	}
	notification.UpdatedAt = now

	// UUID formatında doğrulama
	if _, err := uuid.Parse(notification.UserID); err != nil {
		return fmt.Errorf("invalid UserID format, must be UUID: %w", err)
	}
	if _, err := uuid.Parse(notification.ActorID); err != nil {
		return fmt.Errorf("invalid ActorID format, must be UUID: %w", err)
	}
	if notification.EntityID != "" {
		if _, err := uuid.Parse(notification.EntityID); err != nil {
			return fmt.Errorf("invalid EntityID format, must be UUID: %w", err)
		}
	}

	// MongoDB koleksiyonunu al
	collection := r.GetCollection("notifications")

	// Bildirimi MongoDB'ye ekle
	_, err := collection.InsertOne(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)

	}

	return nil
}

// GetNotificationsByUserID - Belirli bir kullanıcının bildirimlerini getirir
func (r *Repository) GetNotificationsByUserID(ctx context.Context, userID string, limit, skip int64) ([]domain.Notification, error) {
	// UUID doğrulama
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("invalid userID format, must be UUID: %w", err)
	}

	collection := r.GetCollection("notifications")

	// Sorgu oluştur
	filter := bson.M{"userId": userID}

	// Sayfalama ve sıralama seçenekleri
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createdAt": -1}) // En yeni bildirimler önce

	if limit > 0 {
		findOptions.SetLimit(limit)
	}

	if skip > 0 {
		findOptions.SetSkip(skip)
	}

	// Sorguyu çalıştır
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find notifications: %w", err)
	}
	defer cursor.Close(ctx)

	// Sonuçları bir slice'a dönüştür
	var notifications []domain.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, nil
}

// GetUnreadNotificationsByUserID - Belirli bir kullanıcının bildirimlerini getirir
func (r *Repository) GetUnreadNotificationsByUserID(ctx context.Context, userID string, limit, skip int64) ([]domain.Notification, error) {
	// UUID doğrulama
	if _, err := uuid.Parse(userID); err != nil {
		return nil, fmt.Errorf("invalid userID format, must be UUID: %w", err)
	}

	collection := r.GetCollection("notifications")

	// Sorgu oluştur
	filter := bson.M{"userId": userID, "read": false}

	// Sayfalama ve sıralama seçenekleri
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createdAt": -1}) // En yeni bildirimler önce

	if limit > 0 {
		findOptions.SetLimit(limit)
	}

	if skip > 0 {
		findOptions.SetSkip(skip)
	}

	// Sorguyu çalıştır
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find notifications: %w", err)
	}
	defer cursor.Close(ctx)

	// Sonuçları bir slice'a dönüştür
	var notifications []domain.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, nil
}

// MarkNotificationAsRead - Bildirimi okundu olarak işaretler
func (r *Repository) MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error {

	if _, err := uuid.Parse(userID); err != nil {
		return fmt.Errorf("invalid userID format, must be UUID: %w", err)
	}

	collection := r.GetCollection("notifications")

	// Güncelleme sorgusu
	filter := bson.M{"_id": notificationID, "userId": userID}
	update := bson.M{
		"$set": bson.M{
			"read":      true,
			"updatedAt": time.Now(),
		},
	}

	// Güncelleme işlemi
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	// Eğer hiçbir belge etkilenmediyse
	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found with ID: %s", notificationID)
	}

	return nil
}

// DeleteNotification - Bir bildirimi siler
func (r *Repository) DeleteNotification(ctx context.Context, userID, notificationID string) error {
	collection := r.GetCollection("notifications")

	// Silme sorgusu
	filter := bson.M{"_id": notificationID, "userId": userID}

	// Silme işlemi
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	// Eğer hiçbir belge etkilenmediyse
	if result.DeletedCount == 0 {
		return fmt.Errorf("notification not found with ID: %s", notificationID)
	}

	return nil
}

func (r *Repository) ReadAllNotificationsByUserID(ctx context.Context, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return fmt.Errorf("invalid userID FORMAT,MUST BE uuıd: %w ", err)
	}
	collection := r.GetCollection("notifications")
	filter := bson.M{
		"userId": userID,
		"read":   false,
	}
	update := bson.M{
		"$set": bson.M{
			"read":     true,
			"updateAt": time.Now(),
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("faild to update notification: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no notificationsUpdate")
	}
	return nil

}

func (r *Repository) DeleteAllNotificationsByUserID(ctx context.Context, userID string) error {
	collection := r.GetCollection("notifications")

	filter := bson.M{"userId": userID}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}
