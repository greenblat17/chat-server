package chat

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/greenblat17/chat-server/internal/converter"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendMessage send message to chat
func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	err := i.messageService.Send(ctx, converter.ToMessageFromAPI(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}
