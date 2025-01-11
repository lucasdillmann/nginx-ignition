package access_list_repository

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/access_list"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
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
		Where("id = ?", id).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) DeleteByID(id uuid.UUID) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.NewDelete().
		Model((*accessListModel)(nil)).
		Where("id = ?", id).
		Exec(r.ctx)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) FindByName(name string) (*access_list.AccessList, error) {
	var model accessListModel

	err := r.database.Select().
		Model(&model).
		Where("name = ?", name).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) FindPage(pageNumber, pageSize int, searchTerms *string) (*pagination.Page[access_list.AccessList], error) {
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
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []access_list.AccessList
	for _, model := range models {
		result = append(result, *toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, &result), nil
}

func (r *repository) FindAll() (*[]access_list.AccessList, error) {
	var models []accessListModel

	err := r.database.Select().
		Model(&models).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []access_list.AccessList
	for _, model := range models {
		result = append(result, *toDomain(&model))
	}

	return &result, nil
}

func (r *repository) Save(accessList *access_list.AccessList) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.NewMerge().Model(toModel(accessList)).Exec(r.ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
