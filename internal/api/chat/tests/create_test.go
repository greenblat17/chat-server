package tests

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	chatAPI "github.com/greenblat17/chat-server/internal/api/chat"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/service"
	"github.com/greenblat17/chat-server/internal/service/mocks"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	type deps struct {
		chatService    service.ChatService
		messageService service.MessageService
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		usernames       = []string{"username1"}
		chatname        = "chat-name"
		id        int64 = 1

		req = &desc.CreateRequest{
			Usernames: usernames,
			ChatName:  chatname,
		}

		chat = &model.Chat{
			ChatName:  chatname,
			Usernames: usernames,
		}

		resp = &desc.CreateResponse{Id: id}
	)

	tests := []struct {
		name     string
		args     args
		mockFunc func(mc *minimock.Controller) deps
		want     *desc.CreateResponse
		err      error
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: req,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatSrvMock := mocks.NewChatServiceMock(mc)
				chatSrvMock.CreateMock.Expect(ctx, chat).Return(id, nil)

				return deps{
					chatService: chatSrvMock,
				}
			},
			want: resp,
			err:  nil,
		},
		{
			name: "ChatService return error",
			args: args{
				ctx: ctx,
				req: req,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatSrvMock := mocks.NewChatServiceMock(mc)
				chatSrvMock.CreateMock.Expect(ctx, chat).Return(int64(0), assert.AnError)

				return deps{
					chatService: chatSrvMock,
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

			got, err := chatHandler.Create(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
