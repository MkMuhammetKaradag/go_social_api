package chat

import (
	"context"
	"fmt"
	"socialmedia/chat/app/chat/usecase"
	"socialmedia/shared/middlewares"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type ChatWebSocketListenHandler struct {
	usecase usecase.ChatWebSocketListenUseCase
}

type ChatWebSocketListenRequest struct {
}

// type ChatWebSocketListenResponse struct {
// 	Message string `json:"message"`
// }

func NewChatWebSocketListenHandler(usecase usecase.ChatWebSocketListenUseCase) *ChatWebSocketListenHandler {
	return &ChatWebSocketListenHandler{usecase: usecase}
}

func (h *ChatWebSocketListenHandler) Handle(c *websocketFiber.Conn, ctx context.Context, req *ChatWebSocketListenRequest) {

	userID, err := getUserIDFromWS(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	conversationID, err := uuid.Parse(c.Params("chatID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	h.usecase.Execute(c, ctx, userID, conversationID)

}
func getUserIDFromWS(c *websocketFiber.Conn) (uuid.UUID, error) {

	userData, _ := middlewares.GetUserDataFromWS(c)
	return uuid.Parse(userData["id"])
}
