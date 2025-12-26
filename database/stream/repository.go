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

func New(database *database.Database) stream.Repository {
	return &repository{
		database: database,
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

func (r *repository) Save(ctx context.Context, stream *stream.Stream) error {
	exists, err := r.ExistsByID(ctx, stream.ID)
	if err != nil {
		return err
	}

	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	model := toModel(stream)

	if exists {
		_, err = transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(ctx)
	}

	if err != nil {
		return err
	}

	err = r.saveLinkedModels(ctx, transaction, stream)
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
	var models []streamModel

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

	var result []stream.Stream
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
	var models []streamModel

	err := r.database.Select().
		Model(&models).
		Where("enabled = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var result []stream.Stream
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
	count, err := r.database.Select().
		Model((*streamModel)(nil)).
		Where(constants.ByIdFilter, id).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
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

	var routeIDs []uuid.UUID
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

func (r *repository) saveLinkedModels(ctx context.Context, transaction bun.Tx, stream *stream.Stream) error {
	err := r.cleanupLinkedModels(ctx, transaction, stream.ID)
	if err != nil {
		return err
	}

	for _, route := range stream.Routes {
		routeModel := toRouteModel(&route, stream.ID)

		_, err = transaction.NewInsert().Model(routeModel).Exec(ctx)
		if err != nil {
			return err
		}

		for _, backend := range route.Backends {
			backendModel := toBackendModel(&backend, nil, &routeModel.ID)

			_, err = transaction.NewInsert().Model(backendModel).Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	defaultBackendModel := toBackendModel(&stream.DefaultBackend, &stream.ID, nil)
	_, err = transaction.NewInsert().Model(defaultBackendModel).Exec(ctx)

	return err
}

func (r *repository) fillLinkedModels(ctx context.Context, stream *stream.Stream) error {
	var routeModels []streamRouteModel
	err := r.database.Select().
		Model(&routeModels).
		Where(byStreamIdFilter, stream.ID).
		Scan(ctx)
	if err != nil {
		return err
	}

	for _, routeModel := range routeModels {
		var backendModels []streamBackendModel
		err = r.database.
			Select().
			Model(&backendModels).
			Where(byStreamRouteIdFilter, routeModel.ID).
			Scan(ctx)
		if err != nil {
			return err
		}

		stream.Routes = append(stream.Routes, toDomainRoute(&routeModel, backendModels))
	}

	var defaultBackendModel streamBackendModel
	err = r.database.
		Select().
		Model(&defaultBackendModel).
		Where(byStreamIdFilter, stream.ID).
		Scan(ctx)
	if err != nil {
		return err
	}

	stream.DefaultBackend = toDomainBackend(&defaultBackendModel)
	return nil
}
