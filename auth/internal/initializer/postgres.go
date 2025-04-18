package initializer

import (
	"fmt"
	"log"
	"socialmedia/auth/infra/postgres"
	"socialmedia/auth/pkg/config"
)

func InitDatabase(appConfig config.Config) *postgres.Repository {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", appConfig.Database.User, appConfig.Database.Password, appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.DB)
	repo, err := postgres.NewRepository(databaseURL)
	// repo.StartCleanupJob(5 * time.Minute)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	return repo
}
