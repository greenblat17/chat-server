package audit

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/greenblat17/chat-server/internal/model"
	"github.com/greenblat17/chat-server/internal/repository"
	"github.com/greenblat17/platform-common/pkg/db"
)

const (
	auditTable = "audit"

	idColumn        = "id"
	entityColumn    = "entity"
	actionColumn    = "action"
	createdAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates audit repository
func NewRepository(db db.Client) repository.AuditRepository {
	return &repo{db: db}
}

func (r *repo) Save(ctx context.Context, audit *model.Audit) error {
	sqb := sq.Insert(auditTable).
		Columns(entityColumn, actionColumn, createdAtColumn).
		Values(audit.Entity, audit.Action, time.Now()).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sqb.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "AuditRepository.Save",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("failed to save audit")
	}

	return nil
}
