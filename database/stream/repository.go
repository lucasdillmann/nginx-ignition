package stream

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const (
	byStreamIdFilter           = "stream_id = ?"
	byStreamRouteIdFilter      = "stream_route_id = ?"
	byStreamRouteIdArrayFilter = "stream_route_id in (?)"
	byIdArrayFilter            = "id in (?)"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) stream.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*stream.Stream, error) {
	var model streamModel

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

	domain := toDomain(&model)
	err = r.fillLinkedModels(ctx, &domain)
	if err != nil {
		return nil, err
	}

	return &domain, nil
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

	_, err = transaction.
		NewDelete().
		Model((*streamModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(ctx context.Context, strm *stream.Stream) error {
	exists, err := r.ExistsByID(ctx, strm.ID)
	if err != nil {
		return err
	}

	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	model := toModel(strm)

	if exists {
		_, err = transaction.NewUpdate().Model(&model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
	} else {
		_, err = transaction.NewInsert().Model(&model).Exec(ctx)
	}

	if err != nil {
		return err
	}

	err = r.saveLinkedModels(ctx, transaction, strm)
	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindPage(
	ctx context.Context,
	pageSize, pageNumber int,
	searchTerms *string,
) (*pagination.Page[stream.Stream], error) {
	models := make([]streamModel, 0)

	query := r.database.Select().Model(&models)
	if searchTerms != nil && *searchTerms != "" {
		query = query.Where("name ILIKE ?", "%"+*searchTerms+"%")
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

	result := make([]stream.Stream, 0)
	for _, model := range models {
		domainValue := toDomain(&model)
		err = r.fillLinkedModels(ctx, &domainValue)
		if err != nil {
			return nil, err
		}

		result = append(result, domainValue)
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAllEnabled(ctx context.Context) ([]stream.Stream, error) {
	models := make([]streamModel, 0)

	err := r.database.Select().
		Model(&models).
		Where("enabled = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]stream.Stream, 0)
	for _, model := range models {
		domainValue := toDomain(&model)
		err = r.fillLinkedModels(ctx, &domainValue)
		if err != nil {
			return nil, err
		}

		result = append(result, domainValue)
	}

	return result, nil
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.database.Select().
		Model((*streamModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exists(ctx)
}

func (r *repository) cleanupLinkedModels(ctx context.Context, transaction bun.Tx, id uuid.UUID) error {
	_, err := transaction.
		NewDelete().
		Table("stream_backend").
		Where(byStreamIdFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	routeIDs := make([]uuid.UUID, 0)
	err = transaction.
		NewSelect().
		Table("stream_route").
		Column("id").
		Where(byStreamIdFilter, id).
		Scan(ctx, &routeIDs)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if len(routeIDs) != 0 {
		_, err = transaction.
			NewDelete().
			Table("stream_backend").
			Where(byStreamRouteIdArrayFilter, bun.In(routeIDs)).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = transaction.
			NewDelete().
			Table("stream_route").
			Where(byIdArrayFilter, bun.In(routeIDs)).
			Exec(ctx)
	}

	return err
}

func (r *repository) saveLinkedModels(ctx context.Context, transaction bun.Tx, strm *stream.Stream) error {
	err := r.cleanupLinkedModels(ctx, transaction, strm.ID)
	if err != nil {
		return err
	}

	for _, route := range strm.Routes {
		routeModel := toRouteModel(&route, strm.ID)

		_, err = transaction.NewInsert().Model(&routeModel).Exec(ctx)
		if err != nil {
			return err
		}

		for _, backend := range route.Backends {
			backendModel := toBackendModel(&backend, nil, &routeModel.ID)

			_, err = transaction.NewInsert().Model(&backendModel).Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	defaultBackendModel := toBackendModel(&strm.DefaultBackend, &strm.ID, nil)
	_, err = transaction.NewInsert().Model(&defaultBackendModel).Exec(ctx)

	return err
}

func (r *repository) fillLinkedModels(ctx context.Context, strm *stream.Stream) error {
	routeModels := make([]streamRouteModel, 0)
	err := r.database.Select().
		Model(&routeModels).
		Where(byStreamIdFilter, strm.ID).
		Scan(ctx)
	if err != nil {
		return err
	}

	for _, routeModel := range routeModels {
		backendModels := make([]streamBackendModel, 0)
		err = r.database.
			Select().
			Model(&backendModels).
			Where(byStreamRouteIdFilter, routeModel.ID).
			Scan(ctx)
		if err != nil {
			return err
		}

		strm.Routes = append(strm.Routes, toDomainRoute(&routeModel, backendModels))
	}

	var defaultBackendModel streamBackendModel
	err = r.database.
		Select().
		Model(&defaultBackendModel).
		Where(byStreamIdFilter, strm.ID).
		Scan(ctx)
	if err != nil {
		return err
	}

	strm.DefaultBackend = toDomainBackend(&defaultBackendModel)
	return nil
}
