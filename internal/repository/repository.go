package repository

import (
	"context"

	"github.com/greenblat17/chat-server/internal/model"
)

// ChatRepository is a repository for chat
type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// MessageRepository is a repository for messages
type MessageRepository interface {
	Send(ctx context.Context, message *model.Message) error
}

// AuditRepository is a repository for auditing actions
type AuditRepository interface {
	Save(ctx context.Context, audit *model.Audit) error
}

// UserRepository is a repository for users chat
type UserRepository interface {
	SaveByChatID(ctx context.Context, chatID int64, usernames []string) error
	DeleteByChatID(ctx context.Context, chatID int64) error
}
