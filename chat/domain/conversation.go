package domain

import (
	"database/sql"
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
	UserID         uuid.UUID
	SenderUsername string `db:"username"`
	SenderAvatar   sql.NullString `db:"avatar_url"`
	Content        string
	CreatedAt      time.Time
	IsEdited       bool
	DeletedAt      sql.NullTime `db:"deleted_at"`
	ReadAt         sql.NullTime `db:"read_at"`
	Attachments    []Attachment
}
type AttachmentInfo struct {
	ID       uuid.UUID `json:"id"`
	FileURL  string    `json:"file_url"`
	FileType string    `json:"file_type"`
}
type MessageNotification struct {
	Type           string           `json:"type"` // "add", "delete"
	MessageID      uuid.UUID        `json:"message_id"`
	ConversationID uuid.UUID        `json:"conversation_id"`
	UserID         uuid.UUID        `json:"user_id"`
	Username       string           `json:"username,omitempty"`
	Avatar         string           `json:"avatar,omitempty"`
	Content        string           `json:"content,omitempty"`
	CreatedAt      string           `json:"created_at,omitempty"`
	HasAttachments bool             `json:"has_attachments"`
	DeletedAt      string           `json:"deleted_at,omitempty"`
	Attachments    []AttachmentInfo `json:"attachments,omitempty"`
}
type ConversationUserManager struct {
	UserID         uuid.UUID
	ConversationID uuid.UUID
	Username       string
	Avatar         string
	Reason         string
	Type           string
}

type BlockedParticipant struct {
	BlockerID uuid.UUID
	BlockedID uuid.UUID
}
