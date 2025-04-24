package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type chatWebSocketListenUseCase struct {
	repository Repository
	hub        Hub
}

func NewChatWebSocketListenUseCase(repository Repository, hub Hub) ChatWebSocketListenUseCase {
	return &chatWebSocketListenUseCase{
		repository: repository,
		hub:        hub,
	}
}

func (u *chatWebSocketListenUseCase) Execute(c *websocketFiber.Conn, userID, conversationID uuid.UUID) {
	channelName := fmt.Sprintf("conversation:%s", conversationID)
	go u.hub.ListenRedisSendMessage(context.Background(), channelName)

	fmt.Println(userID, conversationID)
	conn := c.Conn
	client := &domain.Client{
		ConversationID: conversationID,
		Conn:           conn,
	}
	fmt.Println(client)
	u.hub.RegisterClient(client)

	defer func() {
		u.hub.UnregisterClient(client)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

}
