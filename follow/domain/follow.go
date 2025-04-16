package domain

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `json:"id"`
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	FollowedAt  time.Time `json:"followed_at"`
}
