package messaging

import (
	"time"
)

type Message struct {
	ID          string      `json:"id"`           // Unique message ID
	Type        MessageType `json:"type"`         // Message type (e.g., "user_created")
	Data        interface{} `json:"data"`         // Actual message payload
	Created     time.Time   `json:"created"`      // Message creation time
	FromService ServiceType `json:"from_service"` // Source service
	// ToService   ServiceType `json:"to_service"`   // Target service (empty for broadcast)
	ToServices []ServiceType `json:"to_services"` // ToService yerine ToServices array'i
	RetryCount int           `json:"retry_count"` // Number of retry attempts
	Priority   int           `json:"priority"`    // Message priority (0-9)
	Headers    Headers       `json:"headers"`     // Custom message headers
	Critical   bool          `json:"critical"`
}

type Headers map[string]interface{}

type MessageHandler func(Message) error
