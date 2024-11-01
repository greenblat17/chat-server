package chat

import (
	"context"
	"errors"

	"github.com/greenblat17/platform-common/pkg/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
)

const (
	chatTable     = "chats"
	chatUserTable = "chat_users"

	idColumn       = "id"
	nameColumn     = "chat_name"
	chatIDColumn   = "chat_id"
	usernameColumn = "username"
)

type repo struct {
	db db.Client
}

// NewRepository creates a chat repository.
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

// Create creates a new chat.
func (r *repo) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	sqb := sq.Insert(chatTable).
		Columns(nameColumn).
		Values(chat.ChatName).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sqb.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "ChatRepository.Create",
		QueryRaw: sql,
	}

	var chatID int64
	err = r.db.DB().ScanOneContext(ctx, &chatID, q, args...)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// Delete deletes a chat.
func (r *repo) Delete(ctx context.Context, id int64) error {
	sqb := sq.Delete(chatUserTable).
		Where(sq.Eq{chatIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sqb.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "ChatRepository.Delete",
		QueryRaw: sql,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("chat not deleted from database")
	}

	return nil
}
