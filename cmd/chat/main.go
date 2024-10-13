package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/greenblat17/chat-server/internal/config"
	"github.com/greenblat17/chat-server/internal/config/env"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/greenblat17/chat-server/internal/repository/chat"
	"github.com/greenblat17/chat-server/internal/repository/message"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository
}

// NewServer creates a new server
func NewServer(chatRepository repository.ChatRepository, messageRepository repository.MessageRepository) *server {
	return &server{
		messageRepository: messageRepository,
		chatRepository:    chatRepository,
	}
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	createdChat := &model.Chat{
		Usernames: req.GetUsernames(),
		Name:      req.GetName(),
	}

	id, err := s.chatRepository.Create(ctx, createdChat)
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := s.chatRepository.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	msg := &model.Message{
		From: req.GetFrom(),
		Text: req.GetText(),
	}

	err := s.messageRepository.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	chatRepo := chat.NewRepository(pool)
	msgRepo := message.NewRepository(pool)
	server := NewServer(chatRepo, msgRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, server)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
