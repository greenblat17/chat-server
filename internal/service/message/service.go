package message

import (
	"github.com/greenblat17/chat-server/internal/repository"
	def "github.com/greenblat17/chat-server/internal/service"
	"github.com/greenblat17/platform-common/pkg/db"
)

var _ def.MessageService = (*service)(nil)

type service struct {
	messageRepository repository.MessageRepository
	auditRepository   repository.AuditRepository
	txManager         db.TxManager
}

// NewService returns a new instance of the message service.
func NewService(
	messageRepository repository.MessageRepository,
	auditRepository repository.AuditRepository,
	txManager db.TxManager,
) *service {
	return &service{
		messageRepository: messageRepository,
		auditRepository:   auditRepository,
		txManager:         txManager,
	}
}
