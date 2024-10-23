package chat

import (
	"context"

	"github.com/greenblat17/chat-server/internal/converter"
	"github.com/greenblat17/chat-server/internal/model"
)

func (s *service) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	var id int64

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.chatRepository.Create(ctx, chat)
		if errTx != nil {
			return errTx
		}

		errTx = s.userRepository.SaveByChatID(ctx, chat.ID, chat.Usernames)
		if errTx != nil {
			return errTx
		}

		errTx = s.auditRepository.Save(ctx, converter.ToAuditFromEntity(model.ChatEntityType, "create"))
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
