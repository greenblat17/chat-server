package message

import (
	"context"

	"github.com/greenblat17/chat-server/internal/converter"
	"github.com/greenblat17/chat-server/internal/model"
)

func (s *service) Send(ctx context.Context, message *model.Message) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.messageRepository.Send(ctx, message)
		if errTx != nil {
			return errTx
		}

		errTx = s.auditRepository.Save(ctx, converter.ToAuditFromEntity(model.ChatEntityType, "send"))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
