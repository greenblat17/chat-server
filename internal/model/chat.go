package model

import "time"

// Chat is a chat with users
type Chat struct {
	ID        int64
	ChatName  string
	Usernames []string
	Messages  []Message
	CreatedAt time.Time
}
