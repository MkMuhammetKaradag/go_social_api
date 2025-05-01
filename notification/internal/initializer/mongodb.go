package initializer

import (
	"fmt"
	"log"
	"socialmedia/notification/infra/mongodb"
	"socialmedia/notification/pkg/config"
)

func InitDatabaseMongo(appConfig config.Config) *mongodb.Repository {

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
		appConfig.Database.User,
		appConfig.Database.Password,
		appConfig.Database.Host,
		appConfig.Database.MongoPort)
	fmt.Println(mongoURI)
	repo, err := mongodb.NewRepository(mongoURI, "notificationdb")

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	return repo
}
