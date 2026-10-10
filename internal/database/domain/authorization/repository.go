package authorization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/database"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) authorization.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindJwtSecret(ctx context.Context) (*string, error) {
	var model configurationModel

	err := r.database.Select().Model(&model).Limit(1).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &model.JwtSecret, nil
}

func (r *repository) SaveJwtSecret(ctx context.Context, secret *string) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	if _, err = transaction.NewTruncateTable().
		Model((*configurationModel)(nil)).
		Exec(ctx); err != nil {
		return err
	}

	model := &configurationModel{
		JwtSecret: *secret,
	}

	if _, err = transaction.NewInsert().Model(model).Exec(ctx); err != nil {
		return err
	}

	return transaction.Commit()
}
