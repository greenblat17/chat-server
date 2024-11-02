package tests

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/greenblat17/chat-server/internal/repository/mocks"
	chatService "github.com/greenblat17/chat-server/internal/service/chat"
	"github.com/greenblat17/platform-common/pkg/db"
	dbMocks "github.com/greenblat17/platform-common/pkg/db/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		chat *model.Chat
	}

	type deps struct {
		chatRepository  repository.ChatRepository
		auditRepository repository.AuditRepository
		userRepository  repository.UserRepository
		txManager       db.TxManager
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		id        int64 = 1
		usernames       = []string{"username1"}

		chat = &model.Chat{
			ID:        id,
			Usernames: usernames,
		}

		audit = &model.Audit{
			Entity: model.ChatEntityType,
			Action: "create",
		}
	)

	tests := []struct {
		name     string
		args     args
		mockFunc func(mc *minimock.Controller) deps
		want     int64
		err      error
	}{
		{
			name: "Success",
			args: args{
				ctx:  ctx,
				chat: chat,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.CreateMock.Expect(ctx, chat).Return(id, nil)

				userRepoMock := mocks.NewUserRepositoryMock(mc)
				userRepoMock.SaveByChatIDMock.Expect(ctx, id, usernames).Return(nil)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)
				auditRepoMock.SaveMock.Expect(ctx, audit).Return(nil)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					chatRepository:  chatRepoMock,
					auditRepository: auditRepoMock,
					userRepository:  userRepoMock,
					txManager:       txManagerMock,
				}
			},
			want: id,
			err:  nil,
		},
		{
			name: "ChatRepository return error",
			args: args{
				ctx:  ctx,
				chat: chat,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.CreateMock.Expect(ctx, chat).Return(int64(0), assert.AnError)

				userRepoMock := mocks.NewUserRepositoryMock(mc)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					chatRepository:  chatRepoMock,
					auditRepository: auditRepoMock,
					userRepository:  userRepoMock,
					txManager:       txManagerMock,
				}
			},
			want: int64(0),
			err:  assert.AnError,
		},
		{
			name: "UserRepository return error",
			args: args{
				ctx:  ctx,
				chat: chat,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.CreateMock.Expect(ctx, chat).Return(id, nil)

				userRepoMock := mocks.NewUserRepositoryMock(mc)
				userRepoMock.SaveByChatIDMock.Expect(ctx, id, usernames).Return(assert.AnError)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					chatRepository:  chatRepoMock,
					auditRepository: auditRepoMock,
					userRepository:  userRepoMock,
					txManager:       txManagerMock,
				}
			},
			want: int64(0),
			err:  assert.AnError,
		},
		{
			name: "AuditRepository return error",
			args: args{
				ctx:  ctx,
				chat: chat,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.CreateMock.Expect(ctx, chat).Return(id, nil)

				userRepoMock := mocks.NewUserRepositoryMock(mc)
				userRepoMock.SaveByChatIDMock.Expect(ctx, id, usernames).Return(nil)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)
				auditRepoMock.SaveMock.Expect(ctx, audit).Return(assert.AnError)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					chatRepository:  chatRepoMock,
					auditRepository: auditRepoMock,
					userRepository:  userRepoMock,
					txManager:       txManagerMock,
				}
			},
			want: int64(0),
			err:  assert.AnError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := tt.mockFunc(mc)

			chatSrv := chatService.NewService(
				deps.userRepository,
				deps.chatRepository,
				deps.auditRepository,
				deps.txManager,
			)

			got, err := chatSrv.Create(tt.args.ctx, tt.args.chat)

			if tt.err != nil {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.want, got)
				require.Nil(t, err)
			}
		})
	}
}
