package settings_repository

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func (r repository) Get() (*settings.Settings, error) {
	transaction, err := r.database.BeginTransaction()
	if err != nil {
		return nil, err
	}

	defer transaction.Rollback()

	var nginxResults, certificateResults, bindingResults, logResults *sql.Rows
	if nginxResults, err = transaction.Query("SELECT * FROM settings_nginx"); err != nil {
		return nil, err
	}

	defer nginxResults.Close()

	if certificateResults, err = transaction.Query("SELECT * FROM settings_certificate_auto_renew"); err != nil {
		return nil, err
	}

	defer certificateResults.Close()

	if bindingResults, err = transaction.Query("SELECT * FROM settings_global_binding"); err != nil {
		return nil, err
	}

	defer bindingResults.Close()

	if logResults, err = transaction.Query("SELECT * FROM settings_log_rotation"); err != nil {
		return nil, err
	}

	defer logResults.Close()

	var nginx *settings.NginxSettings
	if err = nginxResults.Scan(nginx); err != nil {
		return nil, err
	}

	var logRotation *settings.LogRotationSettings
	if err = logResults.Scan(logRotation); err != nil {
		return nil, err
	}

	var certificateAutoRenew *settings.CertificateAutoRenewSettings
	if err = certificateResults.Scan(certificateAutoRenew); err != nil {
		return nil, err
	}

	var globalBindings *[]host.Binding
	for bindingResults.Next() {
		var binding *host.Binding
		if err = bindingResults.Scan(binding); err != nil {
			return nil, err
		}
	}

	result := &settings.Settings{
		Nginx:                nginx,
		LogRotation:          logRotation,
		CertificateAutoRenew: certificateAutoRenew,
		GlobalBindings:       globalBindings,
	}

	return result, nil
}
