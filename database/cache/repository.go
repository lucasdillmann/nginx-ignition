package cache

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const (
	byCacheIDFilter = "cache_id = ?"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) cache.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*cache.Cache, error) {
	var model cacheModel

	err := r.database.Select().
		Model(&model).
		Relation("Durations").
		Where(constants.ByIDFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return ptr.Of(toDomain(&model)), nil
}

func (r *repository) InUseByID(ctx context.Context, id uuid.UUID) (bool, error) {
	hostExists, err := r.database.Select().
		Table("host").
		Where(byCacheIDFilter, id).
		Exists(ctx)
	if err != nil || hostExists {
		return hostExists, err
	}

	return r.database.Select().
		Table("host_route").
		Where(byCacheIDFilter, id).
		Exists(ctx)
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.database.Select().
		Model((*cacheModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exists(ctx)
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	err = r.cleanupLinkedModels(ctx, transaction, id)
	if err != nil {
		return err
	}

	_, err = transaction.NewDelete().
		Model((*cacheModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindPage(
	ctx context.Context,
	pageNumber, pageSize int,
	searchTerms *string,
) (*pagination.Page[cache.Cache], error) {
	models := make([]cacheModel, 0)

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

	result := make([]cache.Cache, 0)
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAllInUse(ctx context.Context) ([]cache.Cache, error) {
	models := make([]cacheModel, 0)

	hostSubquery := r.database.
		Select().
		Table("host").
		Column("cache_id").
		Where("cache_id is not null")
	routeSubquery := r.database.
		Select().
		Table("host_route").
		Column("cache_id").
		Where("cache_id is not null")

	err := r.database.Select().
		Model(&models).
		Relation("Durations").
		Where("id in (?)", hostSubquery).
		WhereOr("id in (?)", routeSubquery).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]cache.Cache, len(models))
	for index, model := range models {
		result[index] = toDomain(&model)
	}

	return result, nil
}

func (r *repository) Save(ctx context.Context, domain *cache.Cache) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	exists, err := transaction.NewSelect().
		Model((*cacheModel)(nil)).
		Where(constants.ByIDFilter, domain.ID).
		Exists(ctx)
	if err != nil {
		return err
	}

	model := toModel(domain)
	if exists {
		err = r.performUpdate(ctx, transaction, &model)
	} else {
		err = r.performInsert(ctx, transaction, &model)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performInsert(
	ctx context.Context,
	transaction bun.Tx,
	model *cacheModel,
) error {
	_, err := transaction.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, transaction, model)
}

func (r *repository) performUpdate(
	ctx context.Context,
	transaction bun.Tx,
	model *cacheModel,
) error {
	_, err := transaction.NewUpdate().
		Model(model).
		Where(constants.ByIDFilter, model.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	err = r.cleanupLinkedModels(ctx, transaction, model.ID)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(ctx, transaction, model)
}

func (r *repository) saveLinkedModels(
	ctx context.Context,
	transaction bun.Tx,
	model *cacheModel,
) error {
	for index := range model.Durations {
		_, err := transaction.NewInsert().
			Model(&model.Durations[index]).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) cleanupLinkedModels(
	ctx context.Context,
	transaction bun.Tx,
	id uuid.UUID,
) error {
	_, err := transaction.NewDelete().
		Table("cache_duration").
		Where(byCacheIDFilter, id).
		Exec(ctx)

	return err
}
