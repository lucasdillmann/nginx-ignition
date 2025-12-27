package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
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
			Read:       &settings.Nginx.Timeouts.Read,
			Connect:    &settings.Nginx.Timeouts.Connect,
			Send:       &settings.Nginx.Timeouts.Send,
			Keepalive:  &settings.Nginx.Timeouts.Keepalive,
			ClientBody: &settings.Nginx.Timeouts.ClientBody,
		},
		Buffers: &nginxBuffersSettingsDto{
			ClientBodyKb:   &settings.Nginx.Buffers.ClientBodyKb,
			ClientHeaderKb: &settings.Nginx.Buffers.ClientHeaderKb,
			LargeClientHeader: &nginxBufferSizeDto{
				SizeKb: &settings.Nginx.Buffers.LargeClientHeader.SizeKb,
				Amount: &settings.Nginx.Buffers.LargeClientHeader.Amount,
			},
			Output: &nginxBufferSizeDto{
				SizeKb: &settings.Nginx.Buffers.Output.SizeKb,
				Amount: &settings.Nginx.Buffers.Output.Amount,
			},
		},
		WorkerProcesses:     &settings.Nginx.WorkerProcesses,
		WorkerConnections:   &settings.Nginx.WorkerConnections,
		DefaultContentType:  &settings.Nginx.DefaultContentType,
		ServerTokensEnabled: &settings.Nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   &settings.Nginx.MaximumBodySizeMb,
		SendfileEnabled:     &settings.Nginx.SendfileEnabled,
		GzipEnabled:         &settings.Nginx.GzipEnabled,
		TcpNoDelayEnabled:   &settings.Nginx.TcpNoDelayEnabled,
		RuntimeUser:         &settings.Nginx.RuntimeUser,
		Custom:              settings.Nginx.Custom,
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

	bindingsModel := make([]bindingDto, 0)
	for _, b := range settings.GlobalBindings {
		bindingsModel = append(bindingsModel, bindingDto{
			Type:          &b.Type,
			IP:            &b.IP,
			Port:          &b.Port,
			CertificateID: b.CertificateID,
		})
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
			Read:       *nginx.Timeouts.Read,
			Connect:    *nginx.Timeouts.Connect,
			Send:       *nginx.Timeouts.Send,
			Keepalive:  *nginx.Timeouts.Keepalive,
			ClientBody: *nginx.Timeouts.ClientBody,
		},
		Buffers: &settings.NginxBuffersSettings{
			ClientBodyKb:   *nginx.Buffers.ClientBodyKb,
			ClientHeaderKb: *nginx.Buffers.ClientHeaderKb,
			LargeClientHeader: &settings.NginxBufferSize{
				SizeKb: *nginx.Buffers.LargeClientHeader.SizeKb,
				Amount: *nginx.Buffers.LargeClientHeader.Amount,
			},
			Output: &settings.NginxBufferSize{
				SizeKb: *nginx.Buffers.Output.SizeKb,
				Amount: *nginx.Buffers.Output.Amount,
			},
		},
		WorkerProcesses:     *nginx.WorkerProcesses,
		WorkerConnections:   *nginx.WorkerConnections,
		DefaultContentType:  *nginx.DefaultContentType,
		ServerTokensEnabled: *nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   *nginx.MaximumBodySizeMb,
		SendfileEnabled:     *nginx.SendfileEnabled,
		GzipEnabled:         *nginx.GzipEnabled,
		TcpNoDelayEnabled:   *nginx.TcpNoDelayEnabled,
		RuntimeUser:         *nginx.RuntimeUser,
		Custom:              nginx.Custom,
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

	globalBindings := make([]binding.Binding, 0)
	for _, b := range bindings {
		globalBindings = append(globalBindings, binding.Binding{
			ID:            uuid.New(),
			Type:          *b.Type,
			IP:            *b.IP,
			Port:          *b.Port,
			CertificateID: b.CertificateID,
		})
	}

	return &settings.Settings{
		Nginx:                nginxSettings,
		LogRotation:          logRotationSettings,
		CertificateAutoRenew: certificateSettings,
		GlobalBindings:       globalBindings,
	}
}
