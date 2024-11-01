package tests

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/greenblat17/chat-server/internal/repository/mocks"
	messageService "github.com/greenblat17/chat-server/internal/service/message"
	"github.com/greenblat17/platform-common/pkg/db"
	dbMocks "github.com/greenblat17/platform-common/pkg/db/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		message *model.Message
	}

	type deps struct {
		messageRepository repository.MessageRepository
		auditRepository   repository.AuditRepository
		txManager         db.TxManager
	}

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		message = &model.Message{
			From: "from",
			Text: "text",
		}
		audit = &model.Audit{
			Entity: model.MessageEntityType,
			Action: "send",
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
				ctx:     ctx,
				message: message,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				messageRepoMock := mocks.NewMessageRepositoryMock(mc)
				messageRepoMock.SendMock.Expect(ctx, message).Return(nil)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)
				auditRepoMock.SaveMock.Expect(ctx, audit).Return(nil)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					messageRepository: messageRepoMock,
					auditRepository:   auditRepoMock,
					txManager:         txManagerMock,
				}
			},
			err: nil,
		},
		{
			name: "MessageRepository return error",
			args: args{
				ctx:     ctx,
				message: message,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				messageRepoMock := mocks.NewMessageRepositoryMock(mc)
				messageRepoMock.SendMock.Expect(ctx, message).Return(assert.AnError)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					messageRepository: messageRepoMock,
					auditRepository:   auditRepoMock,
					txManager:         txManagerMock,
				}
			},
			err: assert.AnError,
		},
		{
			name: "AuditRepository return error",
			args: args{
				ctx:     ctx,
				message: message,
			},
			mockFunc: func(mc *minimock.Controller) deps {
				messageRepoMock := mocks.NewMessageRepositoryMock(mc)
				messageRepoMock.SendMock.Expect(ctx, message).Return(nil)

				auditRepoMock := mocks.NewAuditRepositoryMock(mc)
				auditRepoMock.SaveMock.Expect(ctx, audit).Return(assert.AnError)

				txManagerMock := dbMocks.NewTxManagerMock(mc)
				txManagerMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return deps{
					messageRepository: messageRepoMock,
					auditRepository:   auditRepoMock,
					txManager:         txManagerMock,
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

			chatSrv := messageService.NewService(
				deps.messageRepository,
				deps.auditRepository,
				deps.txManager,
			)

			err := chatSrv.Send(tt.args.ctx, tt.args.message)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
