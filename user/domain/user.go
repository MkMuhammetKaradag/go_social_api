package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username" `
	Email       string    `json:"email" `
	Bio         *string   `json:"bio"`
	AvatarURL   *string   `json:"avatar_url"`
	BannerURL   *string   `json:"banner_url,omitempty"`
	IsPrivate   bool      `json:"is_private" `
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsBlocked   bool
	IsFollowing bool
	Self        bool
}

// Kullanıcı detayları gösterilebilir mi?
func (u *User) CanViewDetails() bool {
	if u.IsBlocked {
		return false
	}

	if !u.IsPrivate {
		return true
	}

	return u.Self || u.IsFollowing
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

type UserSearchResult struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
}
