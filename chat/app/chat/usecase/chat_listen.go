package usecase

import (
	"context"
	"fmt"
	"socialmedia/chat/domain"

	// chatWebsocket "socialmedia/chat/app/chat/websocket"

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
	// Kullanıcının sohbetin bir üyesi olup olmadığını kontrol et
	user, err := u.repository.GetUserIfParticipant(ctx, conversationID, userID)

	if err != nil {
		fmt.Println(err)
		message := websocketFiber.FormatCloseMessage(websocketFiber.CloseInternalServerErr, "Server Error")
		c.Conn.WriteMessage(websocketFiber.CloseMessage, message)
		c.Conn.Close()
		return
	} else if user == nil {
		message := websocketFiber.FormatCloseMessage(websocket.CloseNormalClosure, "Unauthorized access")
		c.Conn.WriteMessage(websocketFiber.CloseMessage, message)
		c.Conn.Close()
		return
	}

	// Conversation üyelerini yükle (eğer daha önce yüklenmemişse)
	if !u.hub.IsConversationLoaded(conversationID) {
		err = u.hub.LoadConversationMembers(ctx, conversationID, u.repository)
		if err != nil {
			message := websocketFiber.FormatCloseMessage(websocket.CloseNormalClosure, "LoadConversationMembers error")
			c.Conn.WriteMessage(websocketFiber.CloseMessage, message)
			c.Conn.Close()
			return
		}
	}

	// Client oluştur ve hub'a kaydet
	conn := c.Conn
	client := &domain.Client{
		ConversationID: conversationID,
		UserID:         userID,
		Avatar:         user.Avatar,
		Username:       user.Username,
		Conn:           conn,
	}

	// Client'i hub'a kaydet (userID ile birlikte)
	u.hub.RegisterClient(client, userID)

	defer func() {
		u.hub.UnregisterClient(client, userID)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
