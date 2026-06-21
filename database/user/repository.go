package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) user.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model userModel

	err := r.database.Select().Model(&model).Where(constants.ByIDFilter, id).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return new(toDomain(&model)), nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*userModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var model userModel

	err := r.database.Select().
		Model(&model).
		Where("username = ?", username).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return new(toDomain(&model)), nil
}

func (r *repository) FindPage(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[user.User], error) {
	models := make([]userModel, 0)

	query := r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where(
			"LOWER(name) LIKE LOWER(?) OR LOWER(username) LIKE LOWER(?)",
			"%"+*searchTerms+"%",
			"%"+*searchTerms+"%",
		)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	query = r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where(
			"LOWER(name) LIKE LOWER(?) OR LOWER(username) LIKE LOWER(?)",
			"%"+*searchTerms+"%",
			"%"+*searchTerms+"%",
		)
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]user.User, 0)
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) IsEnabledByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var model userModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIDFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return model.Enabled, nil
}

func (r *repository) Count(ctx context.Context) (int, error) {
	count, err := r.database.Select().Model((*userModel)(nil)).Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) TryCreateInitialUser(ctx context.Context, u *user.User) (bool, error) {
	transaction, err := r.database.Begin()
	if err != nil {
		return false, err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	if err = lockUserTable(ctx, transaction); err != nil {
		return false, err
	}

	exists, err := transaction.NewSelect().Model((*userModel)(nil)).Exists(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}

	model := toModel(u)
	_, err = transaction.NewInsert().Model(&model).Exec(ctx)
	if err != nil {
		return false, err
	}

	if err := transaction.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func lockUserTable(ctx context.Context, transaction bun.Tx) error {
	var statement string

	switch transaction.Dialect().Name() {
	case dialect.PG:
		statement = `LOCK TABLE "user" IN EXCLUSIVE MODE`
	case dialect.SQLite:
		statement = "ROLLBACK; BEGIN IMMEDIATE;"
	default:
		return errors.New("unsupported database dialect")
	}

	_, err := transaction.ExecContext(ctx, statement)
	return err
}

func (r *repository) Save(ctx context.Context, u *user.User) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	exists, err := transaction.NewSelect().
		Model((*userModel)(nil)).
		Where(constants.ByIDFilter, u.ID).
		Exists(ctx)
	if err != nil {
		return err
	}

	model := toModel(u)
	if exists {
		_, err = transaction.NewUpdate().Model(&model).Where(constants.ByIDFilter, u.ID).Exec(ctx)
	} else {
		_, err = transaction.NewInsert().Model(&model).Exec(ctx)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) TryUpdateLastUsedTOTPCode(
	ctx context.Context,
	id uuid.UUID,
	code string,
) (bool, error) {
	result, err := r.database.Update().
		Model((*userModel)(nil)).
		Set("totp_last_used_codes = SUBSTR(? || COALESCE(',' || totp_last_used_codes, ''), 1, 20)", code).
		Where("id = ? AND (totp_last_used_codes IS NULL OR (',' || totp_last_used_codes || ',') NOT LIKE ?)", id, "%,"+code+",%").
		Exec(ctx)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
