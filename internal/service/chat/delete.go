package chat

import (
	"context"

	"github.com/greenblat17/chat-server/internal/converter"
	"github.com/greenblat17/chat-server/internal/model"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.chatRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.DeleteByChatID(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.auditRepository.Save(ctx, converter.ToAuditFromEntity(model.ChatEntityType, "delete"))
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
