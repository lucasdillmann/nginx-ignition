package user_repository

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/user"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
)

type repository struct {
	database *database.Database
	ctx      context.Context
}

func New(database *database.Database) user.Repository {
	return &repository{
		database: database,
		ctx:      context.Background(),
	}
}

func (r *repository) FindByID(id uuid.UUID) (*user.User, error) {
	var model userModel

	err := r.database.Select().Model(&model).Where("id = ?", id).Scan(r.ctx)
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
		Model((*userModel)(nil)).
		Where("id = ?", id).
		Exec(r.ctx)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) FindByUsername(username string) (*user.User, error) {
	var model userModel

	err := r.database.Select().
		Model(&model).
		Where("username = ?", username).
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	return toDomain(&model), nil
}

func (r *repository) FindPage(pageNumber, pageSize int, searchTerms *string) (*pagination.Page[user.User], error) {
	var models []userModel

	query := r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where("name ILIKE ? OR username ILIKE ?", "%"+*searchTerms+"%", "%"+*searchTerms+"%")
	}

	count, err := query.Count(r.ctx)
	if err != nil {
		return nil, err
	}

	query = r.database.Select().Model(&models)
	if searchTerms != nil {
		query = query.Where("name ILIKE ? OR username ILIKE ?", "%"+*searchTerms+"%", "%"+*searchTerms+"%")
	}

	err = query.
		Limit(pageSize).
		Offset(pageSize * pageNumber).
		Order("name").
		Scan(r.ctx)

	if err != nil {
		return nil, err
	}

	var result []user.User
	for _, model := range models {
		result = append(result, *toDomain(&model))
	}

	return pagination.New(pageNumber, pageSize, count, &result), nil
}

func (r *repository) IsEnabledByID(id uuid.UUID) (bool, error) {
	var model userModel

	err := r.database.Select().
		Model(&model).
		Where("id = ?", id).
		Scan(r.ctx)

	if err != nil {
		return false, err
	}

	return model.Enabled, nil
}

func (r *repository) Count() (int, error) {
	count, err := r.database.Select().Model((*userModel)(nil)).Count(r.ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) Save(user *user.User) error {
	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.NewMerge().Model(toModel(user)).Exec(r.ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
