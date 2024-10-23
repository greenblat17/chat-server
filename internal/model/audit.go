package model

import "time"

// EntityType represents entities that can be existing in database
type EntityType string

var (
	// ChatEntityType represents chat entities
	ChatEntityType EntityType = "chat"
	// MessageEntityType represents message entities
	MessageEntityType EntityType = "message"
)

// Audit is an entity for logging
type Audit struct {
	Entity    EntityType
	Action    string
	CreatedAt time.Time
}
