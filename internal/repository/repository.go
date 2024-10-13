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
