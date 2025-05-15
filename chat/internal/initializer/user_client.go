package initializer

import (
	"fmt"
	"log"
	"socialmedia/chat/internal/grpcclient"
	"socialmedia/chat/pkg/config"
)

func InitUserClient(appConfig config.Config) *grpcclient.UserClient {
	url := fmt.Sprintf("localhost:%s", "8082")
	client, err := grpcclient.NewUserClient(url)

	if err != nil {
		log.Fatalf("user client  failed: %v", err)
	}
	return client
}
