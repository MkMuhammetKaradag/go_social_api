package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/gofiber/contrib/websocket"
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

func (u *chatWebSocketListenUseCase) Execute(c *websocketFiber.Conn, ctx context.Context, userID, conversationID uuid.UUID) {

	isMember, err := u.repository.IsParticipant(ctx, conversationID, userID)

	if err != nil {
		message := websocketFiber.FormatCloseMessage(websocketFiber.CloseInternalServerErr, "Server Error")
		c.Conn.WriteMessage(websocketFiber.CloseMessage, message)
		c.Conn.Close()
		return
	} else if !isMember {
		message := websocketFiber.FormatCloseMessage(websocket.CloseNormalClosure, "Unauthorized access")
		c.Conn.WriteMessage(websocketFiber.CloseMessage, message)
		c.Conn.Close()
		return
	}
	channelName := fmt.Sprintf("conversation:%s", conversationID)
	go u.hub.ListenRedisSendMessage(context.Background(), channelName)
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
