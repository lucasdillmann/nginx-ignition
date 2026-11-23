package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/pagination"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/constants"
	"github.com/lucasdillmann/nginx-ignition/internal/database/core/database"
)

const (
	byUserIDFilter = "user_id = ?"
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
		Model((*sessionModel)(nil)).
		Where(byUserIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*userTokenModel)(nil)).
		Where(byUserIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

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
		Limit(int64(pageSize)).
		Offset(int64(pageSize * pageNumber)).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]user.User, 0)
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, int(count), result), nil
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

	return int(count), nil
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

	if err = r.upsert(ctx, transaction, u); err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) SaveAndDeleteSessions(
	ctx context.Context,
	u *user.User,
	exceptSessionID *uuid.UUID,
) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	if err = r.upsert(ctx, transaction, u); err != nil {
		return err
	}

	query := transaction.NewDelete().
		Model((*sessionModel)(nil)).
		Where(byUserIDFilter, u.ID)
	if exceptSessionID != nil {
		query = query.Where("id <> ?", *exceptSessionID)
	}

	if _, err = query.Exec(ctx); err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) upsert(ctx context.Context, transaction bun.Tx, u *user.User) error {
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

	return err
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

func (r *repository) FindTokensByUserID(
	ctx context.Context,
	userID uuid.UUID,
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[user.APIToken], error) {
	models := make([]userTokenModel, 0)

	query := r.database.Select().
		Model(&models).
		Where(byUserIDFilter, userID)

	if searchTerms != nil {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+*searchTerms+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Limit(int64(pageSize)).
		Offset(int64(pageSize * pageNumber)).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]user.APIToken, 0, len(models))
	for _, model := range models {
		result = append(result, toAPITokenDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, int(count), result), nil
}

func (r *repository) ExistsTokenByName(
	ctx context.Context,
	userID uuid.UUID,
	name string,
) (bool, error) {
	return r.database.Select().
		Model((*userTokenModel)(nil)).
		Where(byUserIDFilter, userID).
		Where("LOWER(name) = LOWER(?)", name).
		Exists(ctx)
}

func (r *repository) FindTokenByID(
	ctx context.Context,
	userID, id uuid.UUID,
) (*user.APIToken, error) {
	var model userTokenModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIDFilter, id).
		Where(byUserIDFilter, userID).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return new(toAPITokenDomain(&model)), nil
}

func (r *repository) CreateToken(ctx context.Context, token *user.APIToken) error {
	model := toAPITokenModel(token)

	_, err := r.database.Insert().Model(&model).Exec(ctx)

	return err
}

func (r *repository) DeleteTokenByID(ctx context.Context, userID, id uuid.UUID) error {
	_, err := r.database.Delete().
		Model((*userTokenModel)(nil)).
		Where(constants.ByIDFilter, id).
		Where(byUserIDFilter, userID).
		Exec(ctx)

	return err
}

func (r *repository) SaveSession(ctx context.Context, session *user.Session) error {
	model := toSessionModel(session)

	_, err := r.database.Insert().Model(&model).Exec(ctx)

	return err
}

func (r *repository) FindSession(ctx context.Context, id uuid.UUID) (*user.Session, error) {
	var model sessionModel

	err := r.database.Select().Model(&model).Where(constants.ByIDFilter, id).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return new(toSessionDomain(&model)), nil
}

func (r *repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	_, err := r.database.Delete().
		Model((*sessionModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exec(ctx)

	return err
}

func (r *repository) RotateSession(
	ctx context.Context,
	oldID uuid.UUID,
	session *user.Session,
) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	result, err := transaction.NewDelete().
		Model((*sessionModel)(nil)).
		Where(constants.ByIDFilter, oldID).
		Where(byUserIDFilter, session.UserID).
		Exec(ctx)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("session no longer exists")
	}

	model := toSessionModel(session)
	if _, err = transaction.NewInsert().Model(&model).Exec(ctx); err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) DeleteExpiredSessions(ctx context.Context) error {
	_, err := r.database.Delete().
		Model((*sessionModel)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)

	return err
}
