package user

import (
	"context"
	"fmt"
	"socialmedia/shared/middlewares"
	"socialmedia/user/app/user/usecase"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type UserStatusPublishHandler struct {
	usecase usecase.UserStatusPublishUseCase
}

type UserStatusPublishRequest struct {
}

func NewUserStatusPublishHandler(usecase usecase.UserStatusPublishUseCase) *UserStatusPublishHandler {
	return &UserStatusPublishHandler{usecase: usecase}
}

func (h *UserStatusPublishHandler) Handle(c *websocketFiber.Conn, ctx context.Context, req *UserStatusPublishRequest) {

	currentUserID, err := getUserIDFromWS(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.usecase.Execute(c, ctx, currentUserID)

}
func getUserIDFromWS(c *websocketFiber.Conn) (uuid.UUID, error) {

	userData, _ := middlewares.GetUserDataFromWS(c)
	return uuid.Parse(userData["id"])
}
