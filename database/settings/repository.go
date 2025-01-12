package settings

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
)

type repository struct {
	database *database.Database
	ctx      context.Context
}

func New(database *database.Database) settings.Repository {
	return &repository{
		database: database,
		ctx:      context.Background(),
	}
}

func (r repository) Get() (*settings.Settings, error) {
	nginx := NginxModel{}
	if err := r.database.Select().Model(&nginx).Scan(r.ctx); err != nil {
		return nil, err
	}

	certificate := CertificateModel{}
	if err := r.database.Select().Model(&certificate).Scan(r.ctx); err != nil {
		return nil, err
	}

	logRotation := LogRotationModel{}
	if err := r.database.Select().Model(&logRotation).Scan(r.ctx); err != nil {
		return nil, err
	}

	var bindings []BindingModel
	if err := r.database.Select().Model(&bindings).Scan(r.ctx); err != nil {
		return nil, err
	}

	for _, binding := range bindings {
		if binding.CertificateID != nil && *binding.CertificateID == uuid.Nil {
			binding.CertificateID = nil
		}
	}

	return toDomain(&nginx, &logRotation, &certificate, &bindings), nil
}

func (r repository) Save(settings *settings.Settings) error {
	nginx, certificate, logRotation, bindings := toModel(settings)

	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	if _, err = transaction.NewTruncateTable().Model(nginx).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewTruncateTable().Model(certificate).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewTruncateTable().Model(logRotation).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewTruncateTable().Model(bindings).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(nginx).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(certificate).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(logRotation).Exec(r.ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(bindings).Exec(r.ctx); err != nil {
		return err
	}

	return transaction.Commit()
}
