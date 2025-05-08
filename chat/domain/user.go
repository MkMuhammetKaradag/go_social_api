package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	UpdatedAt time.Time `json:"updated_at"`
}
