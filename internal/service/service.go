package service

import (
	"context"

	"github.com/greenblat17/chat-server/internal/model"
)

// ChatService is a service for chat
type ChatService interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// MessageService is a service for message
type MessageService interface {
	Send(ctx context.Context, message *model.Message) error
}
