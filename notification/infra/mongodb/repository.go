package mongodb

import (
	"context"
	"fmt"
	"log"

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

// CreateUser fonksiyonu
func (r *Repository) CreateUser(ctx context.Context, id uuid.UUID, username string) error {
	// Kullanıcı koleksiyonuna erişiyoruz
	usersCollection := r.GetCollection("users")
	userId := id.String()
	// Veritabanında kullanıcıyı oluşturma ya da güncelleme işlemi
	filter := bson.M{"id": userId}
	update := bson.M{
		"$set": bson.M{
			"id":        userId,
			"username":  username,
			"updatedAt": "now",
		},
	}
	options := options.Update().SetUpsert(true) // Upsert: varsa güncelle, yoksa ekle

	_, err := usersCollection.UpdateOne(ctx, filter, update, options)
	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}
	return nil
}
