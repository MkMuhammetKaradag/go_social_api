package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	UpdatedAt   time.Time `json:"updated_at"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
}
type FollowRequestUser struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	RequestedAt time.Time `json:"requested_at,omitempty"`
}

