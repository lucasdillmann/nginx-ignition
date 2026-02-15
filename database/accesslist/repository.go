package accesslist

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
)

const (
	byAccessListIDFilter = "access_list_id = ?"
)

type repository struct {
	database *database.Database
}

func New(db *database.Database) accesslist.Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*accesslist.AccessList, error) {
	var model accessListModel

	err := r.database.Select().
		Model(&model).
		Relation("Credentials").
		Relation("EntrySets").
		Where(constants.ByIDFilter, id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return new(toDomain(&model)), nil
}

func (r *repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.database.Select().
		Model((*accessListModel)(nil)).
		Where(constants.ByIDFilter, id).
		Exists(ctx)
}

func (r *repository) InUseByID(ctx context.Context, id uuid.UUID) (bool, error) {
	hostExists, err := r.database.Select().
		Table("host").
		Where(byAccessListIDFilter, id).
		Exists(ctx)
	if err != nil || hostExists {
		return hostExists, err
	}

	return r.database.Select().
		Table("host_route").
		Where(byAccessListIDFilter, id).
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
		Model((*accessListModel)(nil)).
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
) (*pagination.Page[accesslist.AccessList], error) {
	models := make([]accessListModel, 0)

	query := r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+*searchTerms+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Relation("Credentials").
		Relation("EntrySets").
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]accesslist.AccessList, 0)
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAll(ctx context.Context) ([]accesslist.AccessList, error) {
	models := make([]accessListModel, 0)

	err := r.database.Select().
		Model(&models).
		Relation("Credentials").
		Relation("EntrySets").
		Order("name").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]accesslist.AccessList, 0)
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return result, nil
}

func (r *repository) Save(ctx context.Context, accessList *accesslist.AccessList) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer transaction.Rollback()

	exists, err := transaction.NewSelect().
		Model((*accessListModel)(nil)).
		Where(constants.ByIDFilter, accessList.ID).
		Exists(ctx)
	if err != nil {
		return err
	}

	model := toModel(accessList)
	if exists {
		err = r.performUpdate(ctx, &model, transaction)
	} else {
		_, err = transaction.NewInsert().Model(&model).Exec(ctx)
		if err != nil {
			return err
		}

		err = r.saveLinkedModels(ctx, transaction, &model)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performUpdate(
	ctx context.Context,
	model *accessListModel,
	transaction bun.Tx,
) error {
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIDFilter, model.ID).Exec(ctx)
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
	model *accessListModel,
) error {
	for _, credentials := range model.Credentials {
		credentials.ID = uuid.New()
		credentials.AccessListID = model.ID

		_, err := transaction.NewInsert().Model(&credentials).Exec(ctx)
		if err != nil {
			return err
		}
	}

	for _, entrySet := range model.EntrySets {
		entrySet.ID = uuid.New()
		entrySet.AccessListID = model.ID

		_, err := transaction.NewInsert().Model(&entrySet).Exec(ctx)
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
	_, err := transaction.
		NewDelete().
		Table("access_list_credentials").
		Where(byAccessListIDFilter, id).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = transaction.
		NewDelete().
		Table("access_list_entry_set").
		Where(byAccessListIDFilter, id).
		Exec(ctx)

	return err
}
