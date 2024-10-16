package model

import "time"

// Message is a message sent by a user
type Message struct {
	From   string
	Text   string
	SentAt time.Time
}
