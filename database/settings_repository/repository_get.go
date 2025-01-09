package settings_repository

import (
	"database/sql"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"errors"
)

func (r repository) Get() (*settings.Settings, error) {
	transaction, err := r.database.BeginTransaction()
	if err != nil {
		return nil, err
	}

	defer transaction.Rollback()

	nginx, err := fetchNginxSettings(transaction)
	if err != nil {
		return nil, err
	}

	certificate, err := fetchCertificateSettings(transaction)
	if err != nil {
		return nil, err
	}

	bindings, err := fetchBindingSettings(transaction)
	if err != nil {
		return nil, err
	}

	logs, err := fetchLogSettings(transaction)
	if err != nil {
		return nil, err
	}

	result := &settings.Settings{
		Nginx:                nginx,
		LogRotation:          logs,
		CertificateAutoRenew: certificate,
		GlobalBindings:       bindings,
	}

	return result, nil
}

func fetchNginxSettings(transaction *sql.Tx) (*settings.NginxSettings, error) {
	var result *sql.Rows
	var err error

	if result, err = transaction.Query("SELECT * FROM settings_nginx"); err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, errors.New("no results found for nginx settings")
	}

	var nginx *settings.NginxSettings
	if err := result.Scan(nginx); err != nil {
		return nil, err
	}

	return nginx, nil
}

func fetchCertificateSettings(transaction *sql.Tx) (*settings.CertificateAutoRenewSettings, error) {
	var result *sql.Rows
	var err error

	if result, err = transaction.Query("SELECT * FROM settings_certificate_auto_renew"); err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, errors.New("no results found for certificate auto renew settings")
	}

	var certificateAutoRenew *settings.CertificateAutoRenewSettings
	if err = result.Scan(certificateAutoRenew); err != nil {
		return nil, err
	}

	return certificateAutoRenew, nil
}

func fetchBindingSettings(transaction *sql.Tx) (*[]host.Binding, error) {
	var result *sql.Rows
	var err error

	if result, err = transaction.Query("SELECT * FROM settings_global_binding"); err != nil {
		return nil, err
	}

	defer result.Close()

	var globalBindings *[]host.Binding
	for result.Next() {
		var binding *host.Binding
		if err = result.Scan(binding); err != nil {
			return nil, err
		}
	}

	return globalBindings, nil
}

func fetchLogSettings(transaction *sql.Tx) (*settings.LogRotationSettings, error) {
	var result *sql.Rows
	var err error

	if result, err = transaction.Query("SELECT * FROM settings_log_rotation"); err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, errors.New("no results found for log rotation settings")
	}

	var logRotation *settings.LogRotationSettings
	if err = result.Scan(logRotation); err != nil {
		return nil, err
	}

	return logRotation, nil
}
