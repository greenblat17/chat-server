package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler is a function that executed in transaction
type Handler func(ctx context.Context) error

// Client is a client for working with database
type Client interface {
	DB() DB
	Close() error
}

// TxManager a transaction manager that executes a user-specified handler in a transaction
type TxManager interface {
	ReadCommited(ctx context.Context, f Handler) error
}

// Query is a wrapper around a query that stores the query name and the query itself
// The query name is used for logging
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor interface for working with transactions
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer combines NamedExecer and QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer interface for working with named queries using tags in structures
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer interface for working with regular queries
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger interface for checking connection to database
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB is an interface for working with database
type DB interface {
	SQLExecer
	Transactor
	Pinger

	Close()
}
