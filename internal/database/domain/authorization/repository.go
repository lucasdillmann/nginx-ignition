package authorization

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"

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

func (r *repository) SaveJwtSecretIfNotExists(
	ctx context.Context,
	secret *string,
) (*string, error) {
	transaction, err := r.database.Begin()
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	if err = lockAuthorizationConfiguration(ctx, transaction); err != nil {
		return nil, err
	}

	var model configurationModel

	err = transaction.NewSelect().Model(&model).Limit(1).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		model = configurationModel{JwtSecret: *secret}

		if _, err = transaction.NewInsert().Model(&model).Exec(ctx); err != nil {
			return nil, err
		}

		if err = transaction.Commit(); err != nil {
			return nil, err
		}

		return &model.JwtSecret, nil
	}

	if err != nil {
		return nil, err
	}

	return &model.JwtSecret, nil
}

func lockAuthorizationConfiguration(ctx context.Context, transaction bun.Tx) error {
	var statement string

	switch transaction.Dialect().Name() {
	case dialect.PG:
		statement = `LOCK TABLE authorization_configuration IN EXCLUSIVE MODE`
	case dialect.SQLite:
		statement = "ROLLBACK; BEGIN IMMEDIATE;"
	default:
		return errors.New("unsupported database dialect")
	}

	_, err := transaction.ExecContext(ctx, statement)
	return err
}
