package integration

import (
	"context"
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/database/common/constants"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"errors"
)

type repository struct {
	database *database.Database
	ctx      context.Context
}

func New(database *database.Database) integration.Repository {
	return &repository{
		database: database,
		ctx:      context.Background(),
	}
}

func (r *repository) FindByID(id string) (*integration.Integration, error) {
	var model integrationModel

	err := r.database.Select().
		Model(&model).
		Where(constants.ByIdFilter, id).
		Scan(r.ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toDomain(&model)
}

func (r *repository) Save(values *integration.Integration) error {
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	model, err := toModel(values)
	if err != nil {
		return err
	}

	exists, err := transaction.NewSelect().Model((*integrationModel)(nil)).Where(constants.ByIdFilter, values.ID).Exists(r.ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = transaction.NewUpdate().Model(model).Where(constants.ByIdFilter, values.ID).Exec(r.ctx)
	} else {
		_, err = transaction.NewInsert().Model(model).Exec(r.ctx)
	}

	if err != nil {
		return err
	}

	return transaction.Commit()
}
