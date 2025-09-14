package stream

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/stream"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
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

	return toDomain(&model), nil
}

func (r *repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*streamModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(ctx)

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) Save(ctx context.Context, stream *stream.Stream) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	model := toModel(stream)

	exists, err := transaction.NewSelect().Model((*streamModel)(nil)).Where(constants.ByIdFilter, stream.ID).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(ctx)
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
) (*pagination.Page[*stream.Stream], error) {
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

	var result []*stream.Stream
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAllEnabled(ctx context.Context) ([]*stream.Stream, error) {
	var models []streamModel

	err := r.database.Select().
		Model(&models).
		Where("enabled = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	var result []*stream.Stream
	for _, model := range models {
		result = append(result, toDomain(&model))
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
