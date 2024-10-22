package app

import (
	"context"
	"log"

	"github.com/greenblat17/chat-server/internal/api/chat"
	"github.com/greenblat17/chat-server/internal/client/db"
	"github.com/greenblat17/chat-server/internal/client/db/pg"
	"github.com/greenblat17/chat-server/internal/client/db/transaction"
	"github.com/greenblat17/chat-server/internal/closer"
	"github.com/greenblat17/chat-server/internal/config"
	"github.com/greenblat17/chat-server/internal/config/env"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/greenblat17/chat-server/internal/repository/audit"
	chatRepository "github.com/greenblat17/chat-server/internal/repository/chat"
	messageRepository "github.com/greenblat17/chat-server/internal/repository/message"
	userRepository "github.com/greenblat17/chat-server/internal/repository/user"
	"github.com/greenblat17/chat-server/internal/service"
	chatService "github.com/greenblat17/chat-server/internal/service/chat"
	messageService "github.com/greenblat17/chat-server/internal/service/message"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient          db.Client
	txManager         db.TxManager
	userRepository    repository.UserRepository
	auditRepository   repository.AuditRepository
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository

	chatService    service.ChatService
	messageService service.MessageService

	chatImpl *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create database client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) AuditRepository(ctx context.Context) repository.AuditRepository {
	if s.auditRepository == nil {
		s.auditRepository = audit.NewRepository(s.DBClient(ctx))
	}

	return s.auditRepository
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageService == nil {
		s.messageRepository = messageRepository.NewRepository(s.DBClient(ctx))
	}

	return s.messageRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.UserRepository(ctx),
			s.ChatRepository(ctx),
			s.AuditRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) MessageService(ctx context.Context) service.MessageService {
	if s.messageService == nil {
		s.messageService = messageService.NewService(
			s.MessageRepository(ctx),
			s.AuditRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.messageService
}

func (s *serviceProvider) ChatImplementation(ctx context.Context) *chat.Implementation {
	if s.chatImpl == nil {
		s.chatImpl = chat.NewImplementation(
			s.ChatService(ctx),
			s.MessageService(ctx),
		)
	}

	return s.chatImpl
}
