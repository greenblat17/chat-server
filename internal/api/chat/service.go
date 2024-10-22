package chat

import (
	"github.com/greenblat17/chat-server/internal/service"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
)

// Implementation is an implementation of Chat Server.
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService    service.ChatService
	messageService service.MessageService
}

// NewImplementation returns new chat server
func NewImplementation(
	chatService service.ChatService,
	messageService service.MessageService,
) *Implementation {
	return &Implementation{
		chatService:    chatService,
		messageService: messageService,
	}
}
