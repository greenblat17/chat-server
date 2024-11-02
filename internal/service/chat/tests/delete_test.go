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

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
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

		id int64 = 1

		audit = &model.Audit{
			Entity: model.ChatEntityType,
			Action: "delete",
		}
	)

	tests := []struct {
		name     string
		args     args
		mockFunc func(mc *minimock.Controller) deps
		err      error
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.DeleteMock.Expect(ctx, id).Return(nil)

				userRepoMock := mocks.NewUserRepositoryMock(mc)
				userRepoMock.DeleteByChatIDMock.Expect(ctx, id).Return(nil)

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
			err: nil,
		},
		{
			name: "ChatRepository return error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.DeleteMock.Expect(ctx, id).Return(assert.AnError)

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
			err: assert.AnError,
		},
		{
			name: "UserRepository return error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.DeleteMock.Expect(ctx, id).Return(nil)

				userRepoMock := mocks.NewUserRepositoryMock(mc)
				userRepoMock.DeleteByChatIDMock.Expect(ctx, id).Return(assert.AnError)

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
			err: assert.AnError,
		},
		{
			name: "AuditRepository return error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				chatRepoMock := mocks.NewChatRepositoryMock(mc)
				chatRepoMock.DeleteMock.Expect(ctx, id).Return(nil)

				userRepoMock := mocks.NewUserRepositoryMock(mc)
				userRepoMock.DeleteByChatIDMock.Expect(ctx, id).Return(nil)

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
			err: assert.AnError,
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

			err := chatSrv.Delete(tt.args.ctx, tt.args.id)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
