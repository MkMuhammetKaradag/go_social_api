package domain

import (
	"time"
)

type Notification struct {
	ID         string    `bson:"_id,omitempty"`
	UserID     string    `bson:"userId"`     // UUID as string
	Type       string    `bson:"type"`       // like, follow, message, block, mention
	ActorID    string    `bson:"actorId"`    // UUID as string
	EntityID   string    `bson:"entityId"`   // UUID as string (post_id, chat_id vb)
	EntityType string    `bson:"entityType"` // post, chat, comment
	Content    string    `bson:"content"`
	Read       bool      `bson:"read"`
	GroupName  string    `bson:"groupName"`
	Url        string    `bson:"url"`
	CreatedAt  time.Time `bson:"createdAt"`
	UpdatedAt  time.Time `bson:"updatedAt"`
}
