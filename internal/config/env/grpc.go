package env

import (
	"errors"
	"net"
	"os"

	"github.com/greenblat17/chat-server/internal/config"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig creates a new GRPCConfig
func NewGRPCConfig() (*grpcConfig, error) {
	host, ok := os.LookupEnv(grpcHostEnvName)
	if !ok {
		return nil, errors.New("grpc host not found")
	}

	port, ok := os.LookupEnv(grpcPortEnvName)
	if !ok {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
