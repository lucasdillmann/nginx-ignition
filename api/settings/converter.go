package settings

import (
	"github.com/aws/smithy-go/ptr"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func toDto(settings *settings.Settings) *settingsDto {
	if settings == nil {
		return nil
	}

	nginxModel := &nginxSettingsDto{
		Logs: &nginxLogsSettingsDto{
			ServerLogsEnabled: &settings.Nginx.Logs.ServerLogsEnabled,
			ServerLogsLevel:   &settings.Nginx.Logs.ServerLogsLevel,
			AccessLogsEnabled: &settings.Nginx.Logs.AccessLogsEnabled,
			ErrorLogsEnabled:  &settings.Nginx.Logs.ErrorLogsEnabled,
			ErrorLogsLevel:    &settings.Nginx.Logs.ErrorLogsLevel,
		},
		Timeouts: &nginxTimeoutsSettingsDto{
			Read:      &settings.Nginx.Timeouts.Read,
			Connect:   &settings.Nginx.Timeouts.Connect,
			Send:      &settings.Nginx.Timeouts.Send,
			Keepalive: &settings.Nginx.Timeouts.Keepalive,
		},
		WorkerProcesses:     &settings.Nginx.WorkerProcesses,
		WorkerConnections:   &settings.Nginx.WorkerConnections,
		DefaultContentType:  &settings.Nginx.DefaultContentType,
		ServerTokensEnabled: &settings.Nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   &settings.Nginx.MaximumBodySizeMb,
		SendfileEnabled:     &settings.Nginx.SendfileEnabled,
		GzipEnabled:         &settings.Nginx.GzipEnabled,
		RuntimeUser:         ptr.String(string(settings.Nginx.RuntimeUser)),
	}

	logRotationModel := &logRotationSettingsDto{
		Enabled:           &settings.LogRotation.Enabled,
		MaximumLines:      &settings.LogRotation.MaximumLines,
		IntervalUnit:      &settings.LogRotation.IntervalUnit,
		IntervalUnitCount: &settings.LogRotation.IntervalUnitCount,
	}

	certificateModel := &certificateAutoRenewSettingsDto{
		Enabled:           &settings.CertificateAutoRenew.Enabled,
		IntervalUnit:      &settings.CertificateAutoRenew.IntervalUnit,
		IntervalUnitCount: &settings.CertificateAutoRenew.IntervalUnitCount,
	}

	var bindingsModel []*bindingDto
	if settings.GlobalBindings != nil {
		for _, binding := range settings.GlobalBindings {
			bindingsModel = append(bindingsModel, &bindingDto{
				Type:          &binding.Type,
				IP:            &binding.IP,
				Port:          &binding.Port,
				CertificateID: binding.CertificateID,
			})
		}
	}

	return &settingsDto{
		Nginx:                nginxModel,
		LogRotation:          logRotationModel,
		CertificateAutoRenew: certificateModel,
		GlobalBindings:       bindingsModel,
	}
}

func toDomain(input *settingsDto) *settings.Settings {
	nginx := input.Nginx
	logRotation := input.LogRotation
	certificate := input.CertificateAutoRenew
	bindings := input.GlobalBindings

	if nginx == nil || logRotation == nil || certificate == nil {
		return nil
	}

	nginxSettings := &settings.NginxSettings{
		Logs: &settings.NginxLogsSettings{
			ServerLogsEnabled: *nginx.Logs.ServerLogsEnabled,
			ServerLogsLevel:   *nginx.Logs.ServerLogsLevel,
			AccessLogsEnabled: *nginx.Logs.AccessLogsEnabled,
			ErrorLogsEnabled:  *nginx.Logs.ErrorLogsEnabled,
			ErrorLogsLevel:    *nginx.Logs.ErrorLogsLevel,
		},
		Timeouts: &settings.NginxTimeoutsSettings{
			Read:      *nginx.Timeouts.Read,
			Connect:   *nginx.Timeouts.Connect,
			Send:      *nginx.Timeouts.Send,
			Keepalive: *nginx.Timeouts.Keepalive,
		},
		WorkerProcesses:     *nginx.WorkerProcesses,
		WorkerConnections:   *nginx.WorkerConnections,
		DefaultContentType:  *nginx.DefaultContentType,
		ServerTokensEnabled: *nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   *nginx.MaximumBodySizeMb,
		SendfileEnabled:     *nginx.SendfileEnabled,
		GzipEnabled:         *nginx.GzipEnabled,
		RuntimeUser:         settings.RuntimeUser(*nginx.RuntimeUser),
	}

	logRotationSettings := &settings.LogRotationSettings{
		Enabled:           *logRotation.Enabled,
		MaximumLines:      *logRotation.MaximumLines,
		IntervalUnit:      *logRotation.IntervalUnit,
		IntervalUnitCount: *logRotation.IntervalUnitCount,
	}

	certificateSettings := &settings.CertificateAutoRenewSettings{
		Enabled:           *certificate.Enabled,
		IntervalUnit:      *certificate.IntervalUnit,
		IntervalUnitCount: *certificate.IntervalUnitCount,
	}

	var globalBindings []*host.Binding
	if bindings != nil {
		for _, binding := range bindings {
			globalBindings = append(globalBindings, &host.Binding{
				ID:            uuid.New(),
				Type:          *binding.Type,
				IP:            *binding.IP,
				Port:          *binding.Port,
				CertificateID: binding.CertificateID,
			})
		}
	}

	return &settings.Settings{
		Nginx:                nginxSettings,
		LogRotation:          logRotationSettings,
		CertificateAutoRenew: certificateSettings,
		GlobalBindings:       globalBindings,
	}
}
