package message

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	messageTable = "messages"

	idColumn     = "id"
	textColumn   = "text"
	fromColumn   = "from_username"
	sentAtColumn = "sent_at"
)

type repo struct {
	db *pgxpool.Pool
}

// NewRepository creates a new user repository.
func NewRepository(db *pgxpool.Pool) repository.MessageRepository {
	return &repo{db: db}
}

// Send sends a message in the chat
func (r *repo) Send(ctx context.Context, message *model.Message) error {
	sqb := sq.Insert(messageTable).
		Columns(textColumn, fromColumn).
		Values(message.Text, message.From).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sqb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
