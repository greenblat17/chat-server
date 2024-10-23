package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/greenblat17/chat-server/internal/client/db"
)

const (
	chatUserTable = "chat_users"

	chatIDColumn   = "chat_id"
	usernameColumn = "username"
)

type repo struct {
	db db.Client
}

// NewRepository creates a user repository
func NewRepository(db db.Client) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) SaveByChatID(ctx context.Context, chatID int64, usernames []string) error {
	sqb := sq.Insert(chatUserTable).
		Columns(chatIDColumn, usernameColumn).
		PlaceholderFormat(sq.Dollar)

	for _, username := range usernames {
		sqb = sqb.Values(chatID, username)
	}

	sql, args, err := sqb.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "UserRepository.SaveByChatID",
		QueryRaw: sql,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("failed to save usernames for chat")
	}

	return nil
}

func (r *repo) DeleteByChatID(ctx context.Context, chatID int64) error {
	sqb := sq.Delete(chatUserTable).
		Where(sq.Eq{chatIDColumn: chatID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sqb.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "UserRepository.DeleteByChatID",
		QueryRaw: sql,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("failed to delete usernames for chat")
	}

	return nil
}
