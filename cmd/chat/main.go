package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50052
)

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, _ *desc.CreateRequest) (*desc.CreateResponse, error) {
	return &desc.CreateResponse{
		Id: 1,
	}, nil
}

func (s *server) Delete(_ context.Context, _ *desc.DeleteRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, _ *desc.SendMessageRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
