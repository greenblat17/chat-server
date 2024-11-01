package chat

import (
	"github.com/greenblat17/chat-server/internal/repository"
	def "github.com/greenblat17/chat-server/internal/service"
	"github.com/greenblat17/platform-common/pkg/db"
)

var _ def.ChatService = (*service)(nil)

type service struct {
	chatRepository  repository.ChatRepository
	auditRepository repository.AuditRepository
	userRepository  repository.UserRepository
	txManager       db.TxManager
}

// NewService returns a new instance of chat service
func NewService(
	userRepository repository.UserRepository,
	chatRepository repository.ChatRepository,
	auditRepository repository.AuditRepository,
	txManager db.TxManager,
) *service {
	return &service{
		userRepository:  userRepository,
		chatRepository:  chatRepository,
		auditRepository: auditRepository,
		txManager:       txManager,
	}
}
