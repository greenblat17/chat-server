package message

import (
	"context"
	"errors"

	"github.com/greenblat17/chat-server/internal/client/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
)

const (
	messageTable = "messages"

	idColumn     = "id"
	textColumn   = "text"
	fromColumn   = "from_username"
	sentAtColumn = "sent_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates a user repository.
func NewRepository(db db.Client) repository.MessageRepository {
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

	q := db.Query{
		Name:     "MessageRepository.Send",
		QueryRaw: sql,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("failed to save in db sending message")
	}

	return nil
}
