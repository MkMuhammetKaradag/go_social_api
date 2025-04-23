package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID `json:"id"`
	IsGroup   bool      `json:"is_group"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
