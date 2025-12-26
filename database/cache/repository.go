package cache

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) cache.Repository {
	return &repository{
		database: database,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*cache.Cache, error) {
	var model cacheModel

	err := r.database.Select().
		Model(&model).
		Relation("Durations").
		Where(constants.ByIdFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) InUseByID(ctx context.Context, id uuid.UUID) (bool, error) {
	hostExists, err := r.database.Select().
		Table("host").
		Where("cache_id = ?", id).
		Exists(ctx)
	if err != nil || hostExists {
		return hostExists, err
	}

	return r.database.Select().
		Table("host_route").
		Where("cache_id = ?", id).
		Exists(ctx)
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	count, err := r.database.Select().
		Model((*cacheModel)(nil)).
		Where(constants.ByIdFilter, id).
		Count(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	err = r.cleanupLinkedModels(ctx, transaction, id)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*cacheModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindByName(ctx context.Context, name string) (*cache.Cache, error) {
	var model cacheModel

	err := r.database.Select().
		Model(&model).
		Relation("Durations").
		Where("name = ?", name).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) FindPage(
	ctx context.Context,
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[*cache.Cache], error) {
	var models []cacheModel

	query := r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where("name ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Relation("Durations").
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var result []*cache.Cache
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAll(ctx context.Context) ([]*cache.Cache, error) {
	var models []cacheModel

	err := r.database.Select().
		Model(&models).
		Relation("Durations").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var result []*cache.Cache
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return result, nil
}

func (r *repository) Save(ctx context.Context, domain *cache.Cache) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	exists, err := transaction.NewSelect().Model((*cacheModel)(nil)).Where(constants.ByIdFilter, domain.ID).Exists(ctx)
	if err != nil {
		return err
	}

	model := toModel(domain)
	if exists {
		err = r.performUpdate(ctx, transaction, model)
	} else {
		err = r.performInsert(ctx, transaction, model)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performInsert(ctx context.Context, transaction bun.Tx, model *cacheModel) error {
	_, err := transaction.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, transaction, model)
}

func (r *repository) performUpdate(ctx context.Context, transaction bun.Tx, model *cacheModel) error {
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
	if err != nil {
		return err
	}

	err = r.cleanupLinkedModels(ctx, transaction, model.ID)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, transaction, model)
}

func (r *repository) saveLinkedModels(ctx context.Context, transaction bun.Tx, model *cacheModel) error {
	for _, duration := range model.Durations {
		_, err := transaction.NewInsert().Model(duration).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) cleanupLinkedModels(ctx context.Context, transaction bun.Tx, id uuid.UUID) error {
	_, err := transaction.
		NewDelete().
		Table("cache_duration").
		Where("cache_id = ?", id).
		Exec(ctx)

	return err
}
