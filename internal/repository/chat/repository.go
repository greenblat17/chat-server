package chat

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	chatTable     = "chats"
	chatUserTable = "chat_users"

	idColumn       = "id"
	nameColumn     = "name"
	chatIDColumn   = "chat_id"
	usernameColumn = "username"
)

type repo struct {
	db *pgxpool.Pool
}

// NewRepository creates a new user repository.
func NewRepository(db *pgxpool.Pool) repository.ChatRepository {
	return &repo{db: db}
}

// Create creates a new chat.
func (r *repo) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	// Начало транзакции
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// SQL для вставки чата
	sqb := sq.Insert(chatTable).
		Columns(nameColumn).
		Values(chat.Name).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sqb.ToSql()
	if err != nil {
		return 0, err
	}

	// Выполнение запроса на вставку чата
	var chatID int64
	err = tx.QueryRow(ctx, sql, args...).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	// Вставка пользователей чата
	for _, username := range chat.Usernames {
		sqb = sq.Insert(chatUserTable).
			Columns(chatIDColumn, usernameColumn).
			Values(chatID, username).
			PlaceholderFormat(sq.Dollar)

		sql, args, err = sqb.ToSql()
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return 0, err
		}
	}

	// Если всё прошло успешно, фиксируем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// Delete deletes a chat.
func (r *repo) Delete(ctx context.Context, id int64) error {
	// Начало транзакции
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Удаление пользователей, связанных с чатом
	sqb := sq.Delete(chatUserTable).
		Where(sq.Eq{chatIDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := sqb.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	// Удаление самого чата
	sqb = sq.Delete(chatTable).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err = sqb.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	// Если всё прошло успешно, фиксируем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
