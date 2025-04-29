package domain

import (
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
)

type Client struct {
	UserID    uuid.UUID
	Conn      *websocket.Conn
	WriteLock sync.Mutex
}
