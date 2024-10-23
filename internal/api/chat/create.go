package chat

import (
	"context"

	"github.com/greenblat17/chat-server/internal/converter"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create creates a new chat.
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, converter.ToChatFromCreateAPI(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
