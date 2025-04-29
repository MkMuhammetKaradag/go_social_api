package usecase

import (
	"context"
	"socialmedia/user/domain"

	websocketFiber "github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type userStatusPublishUseCase struct {
	repository Repository
	hub        Hub
}

func NewUserStatusPublishUseCase(repository Repository, hub Hub) UserStatusPublishUseCase {
	return &userStatusPublishUseCase{
		repository: repository,
		hub:        hub,
	}
}

func (u *userStatusPublishUseCase) Execute(c *websocketFiber.Conn, ctx context.Context, currentUserID uuid.UUID) {

	conn := c.Conn
	client := &domain.Client{
		UserID: currentUserID,
		Conn:   conn,
	}
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
