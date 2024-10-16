package env

import (
	"errors"
	"os"

	"github.com/greenblat17/chat-server/internal/config"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig creates a new PGConfig
func NewPGConfig() (*pgConfig, error) {
	dsn, ok := os.LookupEnv(dsnEnvName)
	if !ok {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
