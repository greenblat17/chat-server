package tests

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	chatAPI "github.com/greenblat17/chat-server/internal/api/chat"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/service"
	"github.com/greenblat17/chat-server/internal/service/mocks"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	type deps struct {
		chatService    service.ChatService
		messageService service.MessageService
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		fromObject = "from"
		text       = "text"

		req = &desc.SendMessageRequest{
			From: fromObject,
			Text: text,
		}
		message = &model.Message{
			From: fromObject,
			Text: text,
		}
	)

	tests := []struct {
		name     string
		args     args
		mockFunc func(mc *minimock.Controller) deps
		want     *empty.Empty
		err      error
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: req,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				messageSrvMock := mocks.NewMessageServiceMock(mc)
				messageSrvMock.SendMock.Expect(ctx, message).Return(nil)

				return deps{
					messageService: messageSrvMock,
				}
			},
			want: &empty.Empty{},
			err:  nil,
		},
		{
			name: "MessageService return error",
			args: args{
				ctx: ctx,
				req: req,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				messageSrvMock := mocks.NewMessageServiceMock(mc)
				messageSrvMock.SendMock.Expect(ctx, message).Return(assert.AnError)

				return deps{
					messageService: messageSrvMock,
				}
			},
			want: nil,
			err:  status.Error(codes.Internal, assert.AnError.Error()),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := tt.mockFunc(mc)

			chatHandler := chatAPI.NewImplementation(deps.chatService, deps.messageService)

			got, err := chatHandler.SendMessage(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
