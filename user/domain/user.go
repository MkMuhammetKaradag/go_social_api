package domain

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string         `json:"id"`
	Username  string         `json:"username" `
	Email     string         `json:"email" `
	Bio       sql.NullString `json:"bio"`
	AvatarURL sql.NullString `json:"avatar_url"`
	IsPrivate bool           `json:"is_private" `
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
type UserUpdate struct {
	// ID        uuid.UUID
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Location  *string `json:"location,omitempty"`
	Website   *string `json:"website,omitempty"`
	IsPrivate *bool   `json:"is_private,omitempty"`
}
