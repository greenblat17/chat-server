package config

import "github.com/joho/godotenv"

// Load loads configuration from environment variables file
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig is the configuration for the gRPC server
type GRPCConfig interface {
	Address() string
}

// PGConfig is the configuration for the PostgreSQL database
type PGConfig interface {
	DSN() string
}
