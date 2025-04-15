package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" `
	Email     string    `json:"email" `
	Bio       *string   `json:"bio"`
	AvatarURL *string   `json:"avatar_url"`
	BannerURL *string   `json:"banner_url,omitempty"`
	IsPrivate bool      `json:"is_private" `
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type UserUpdate struct {
	// ID        uuid.UUID   banner_url
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	BannerURL *string `json:"banner_url,omitempty"`
	Location  *string `json:"location,omitempty"`
	Website   *string `json:"website,omitempty"`
	IsPrivate *bool   `json:"is_private,omitempty"`
}
