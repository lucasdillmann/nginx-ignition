package access_list

import (
	"context"
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type repository struct {
	database *database.Database
	ctx      context.Context
}

func New(database *database.Database) access_list.Repository {
	return &repository{
		database: database,
		ctx:      context.Background(),
	}
}

func (r *repository) FindByID(id uuid.UUID) (*access_list.AccessList, error) {
	var model accessListModel

	err := r.database.Select().
		Model(&model).
		Relation("Credentials").
		Relation("EntrySets").
		Where(constants.ByIdFilter, id).
		Scan(r.ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) DeleteByID(id uuid.UUID) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	_, err = transaction.NewDelete().
		Model((*accessListModel)(nil)).
		Where(constants.ByIdFilter, id).
		Exec(r.ctx)

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) FindByName(name string) (*access_list.AccessList, error) {
	var model accessListModel

	err := r.database.Select().
		Model(&model).
		Relation("Credentials").
		Relation("EntrySets").
		Where("name = ?", name).
		Scan(r.ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) FindPage(pageNumber, pageSize int, searchTerms *string) (*pagination.Page[*access_list.AccessList], error) {
	var models []accessListModel

	query := r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where("name ILIKE ?", "%"+*searchTerms+"%")
	}

	count, err := query.Count(r.ctx)
	if err != nil {
		return nil, err
	}

	err = query.
		Relation("Credentials").
		Relation("EntrySets").
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []*access_list.AccessList
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, result), nil
}

func (r *repository) FindAll() ([]*access_list.AccessList, error) {
	var models []accessListModel

	err := r.database.Select().
		Model(&models).
		Relation("Credentials").
		Relation("EntrySets").
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []*access_list.AccessList
	for _, model := range models {
		result = append(result, toDomain(&model))
	}

	return result, nil
}

func (r *repository) Save(accessList *access_list.AccessList) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	exists, err := transaction.NewSelect().Model((*accessListModel)(nil)).Where(constants.ByIdFilter, accessList.ID).Exists(r.ctx)
	if err != nil {
		return err
	}

	if exists {
		err = r.performUpdate(toModel(accessList), transaction)
	} else {
		model := toModel(accessList)
		_, err = transaction.NewInsert().Model(model).Exec(r.ctx)
		if err != nil {
			return err
		}

		err = r.saveLinkedModels(transaction, model)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}

func (r *repository) performUpdate(model *accessListModel, transaction bun.Tx) error {
	_, err := transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, model.ID).Exec(r.ctx)
	if err != nil {
		return err
	}

	_, err = transaction.
		NewDelete().
		Table("access_list_credentials").
		Where("access_list_id = ?", model.ID).
		Exec(r.ctx)
	if err != nil {
		return err
	}

	_, err = transaction.
		NewDelete().
		Table("access_list_entry_set").
		Where("access_list_id = ?", model.ID).
		Exec(r.ctx)
	if err != nil {
		return err
	}

	return r.saveLinkedModels(transaction, model)
}

func (r *repository) saveLinkedModels(transaction bun.Tx, model *accessListModel) error {
	for _, credentials := range model.Credentials {
		credentials.ID = uuid.New()
		credentials.AccessListID = model.ID

		_, err := transaction.NewInsert().Model(credentials).Exec(r.ctx)
		if err != nil {
			return err
		}
	}

	for _, entrySet := range model.EntrySets {
		entrySet.ID = uuid.New()
		entrySet.AccessListID = model.ID

		_, err := transaction.NewInsert().Model(entrySet).Exec(r.ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
