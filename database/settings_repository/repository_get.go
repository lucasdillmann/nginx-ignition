package settings_repository

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/settings"
	"github.com/google/uuid"
)

func (r repository) Get() (*settings.Settings, error) {
	ctx := context.Background()

	nginx := NginxModel{}
	if err := r.database.Select().Model(&nginx).Scan(ctx); err != nil {
		return nil, err
	}

	certificate := CertificateModel{}
	if err := r.database.Select().Model(&certificate).Scan(ctx); err != nil {
		return nil, err
	}

	logRotation := LogRotationModel{}
	if err := r.database.Select().Model(&logRotation).Scan(ctx); err != nil {
		return nil, err
	}

	var bindings []BindingModel
	if err := r.database.Select().Model(&bindings).Scan(ctx); err != nil {
		return nil, err
	}

	for _, binding := range bindings {
		if binding.CertificateID != nil && *binding.CertificateID == uuid.Nil {
			binding.CertificateID = nil
		}
	}

	return toDomain(&nginx, &logRotation, &certificate, &bindings), nil
}
