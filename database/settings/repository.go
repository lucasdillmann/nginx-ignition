package settings

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/settings"
	"dillmann.com.br/nginx-ignition/database/common/database"
	"github.com/google/uuid"
)

type repository struct {
	database *database.Database
}

func New(database *database.Database) settings.Repository {
	return &repository{
		database: database,
	}
}

func (r repository) Get(ctx context.Context) (*settings.Settings, error) {
	nginx := nginxModel{}
	if err := r.database.Select().Model(&nginx).Scan(ctx); err != nil {
		return nil, err
	}

	certificate := certificateModel{}
	if err := r.database.Select().Model(&certificate).Scan(ctx); err != nil {
		return nil, err
	}

	logRotation := logRotationModel{}
	if err := r.database.Select().Model(&logRotation).Scan(ctx); err != nil {
		return nil, err
	}

	var bindings []*bindingModel
	if err := r.database.Select().Model(&bindings).Scan(ctx); err != nil {
		return nil, err
	}

	for _, binding := range bindings {
		if binding.CertificateID != nil && *binding.CertificateID == uuid.Nil {
			binding.CertificateID = nil
		}
	}

	return toDomain(&nginx, &logRotation, &certificate, bindings), nil
}

func (r repository) Save(ctx context.Context, settings *settings.Settings) error {
	nginx, certificate, logRotation, bindings := toModel(settings)

	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer transaction.Rollback()

	if _, err = transaction.NewTruncateTable().Model(nginx).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewTruncateTable().Model(certificate).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewTruncateTable().Model(logRotation).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewTruncateTable().Model(&bindings).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(nginx).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(certificate).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(logRotation).Exec(ctx); err != nil {
		return err
	}

	if _, err = transaction.NewInsert().Model(&bindings).Exec(ctx); err != nil {
		return err
	}

	return transaction.Commit()
}
