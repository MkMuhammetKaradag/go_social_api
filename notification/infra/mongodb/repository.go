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
