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

type Attachment struct {
	ID        uuid.UUID
	MessageID uuid.UUID
	FileURL   string
	FileType  string
}

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	CreatedAt      time.Time
	IsEdited       bool
	DeletedAt      *time.Time
	Attachments    []Attachment
}

type BlockedParticipant struct {
	BlockerID uuid.UUID
	BlockedID uuid.UUID
}
