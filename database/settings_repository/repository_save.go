package settings_repository

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func (r repository) Save(settings *settings.Settings) error {
	ctx := context.Background()
	nginx, certificate, logRotation, bindings := toModel(settings)

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.NewTruncateTable().Model(nginx).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewTruncateTable().Model(certificate).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewTruncateTable().Model(logRotation).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewTruncateTable().Model(bindings).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(nginx).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(certificate).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(logRotation).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(bindings).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
