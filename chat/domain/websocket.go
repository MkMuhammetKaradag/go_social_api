package domain

import (
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
)

type Client struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Username       string
	Avatar         string
	Conn           *websocket.Conn
	WriteLock      sync.Mutex
}
