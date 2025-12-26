package integration

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) integration.Repository {
	return &repository{
		database: database,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*integration.Integration, error) {
	var model integrationModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIdFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) ExistsByName(ctx context.Context, name string) (*bool, error) {
	count, err := r.database.Select().
		Model((*integrationModel)(nil)).
		Where("name = ?", name).
		Count(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return ptr.Of(false), nil
	}

	if err != nil {
		return nil, err
	}

	return ptr.Of(count > 0), nil
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (*bool, error) {
	count, err := r.database.Select().
		Model((*integrationModel)(nil)).
		Where(constants.ByIdFilter, id).
		Count(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return ptr.Of(false), nil
	}

	if err != nil {
		return nil, err
	}

	return ptr.Of(count > 0), nil
}

func (r *repository) InUseByID(ctx context.Context, id uuid.UUID) (*bool, error) {
	count, err := r.database.Select().
		Table("host_route").
		Where("integration_id = ? and type = ?", id, "INTEGRATION").
		Count(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return ptr.Of(false), nil
	}

	if err != nil {
		return nil, err
	}

	return ptr.Of(count > 0), nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer tx.Rollback()

	_, err = tx.NewDelete().Model((*integrationModel)(nil)).Where(constants.ByIdFilter, id).Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) Save(ctx context.Context, values *integration.Integration) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model, err := toModel(values)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*integrationModel)(nil)).Where(constants.ByIdFilter, values.ID).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, values.ID).Exec(ctx)
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(ctx)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindPage(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
	enabledOnly bool,
) (*pagination.Page[integration.Integration], error) {
	models := make([]integrationModel, 0)

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("name ILIKE ?", "%"+*searchTerms+"%")
	}

	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]integration.Integration, 0)
	for _, model := range models {
		domain, err := toDomain(&model)
		if err != nil {
			return nil, err
		}

		result = append(result, *domain)
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}
