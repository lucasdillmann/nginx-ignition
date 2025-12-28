package settings

import (
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func toDto(set *settings.Settings) *settingsDto {
	if set == nil {
		return nil
	}

	nginxModel := &nginxSettingsDto{
		Logs: &nginxLogsSettingsDto{
			ServerLogsEnabled: &set.Nginx.Logs.ServerLogsEnabled,
			ServerLogsLevel:   &set.Nginx.Logs.ServerLogsLevel,
			AccessLogsEnabled: &set.Nginx.Logs.AccessLogsEnabled,
			ErrorLogsEnabled:  &set.Nginx.Logs.ErrorLogsEnabled,
			ErrorLogsLevel:    &set.Nginx.Logs.ErrorLogsLevel,
		},
		Timeouts: &nginxTimeoutsSettingsDto{
			Read:       &set.Nginx.Timeouts.Read,
			Connect:    &set.Nginx.Timeouts.Connect,
			Send:       &set.Nginx.Timeouts.Send,
			Keepalive:  &set.Nginx.Timeouts.Keepalive,
			ClientBody: &set.Nginx.Timeouts.ClientBody,
		},
		Buffers: &nginxBuffersSettingsDto{
			ClientBodyKb:   &set.Nginx.Buffers.ClientBodyKb,
			ClientHeaderKb: &set.Nginx.Buffers.ClientHeaderKb,
			LargeClientHeader: &nginxBufferSizeDto{
				SizeKb: &set.Nginx.Buffers.LargeClientHeader.SizeKb,
				Amount: &set.Nginx.Buffers.LargeClientHeader.Amount,
			},
			Output: &nginxBufferSizeDto{
				SizeKb: &set.Nginx.Buffers.Output.SizeKb,
				Amount: &set.Nginx.Buffers.Output.Amount,
			},
		},
		WorkerProcesses:     &set.Nginx.WorkerProcesses,
		WorkerConnections:   &set.Nginx.WorkerConnections,
		DefaultContentType:  &set.Nginx.DefaultContentType,
		ServerTokensEnabled: &set.Nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   &set.Nginx.MaximumBodySizeMb,
		SendfileEnabled:     &set.Nginx.SendfileEnabled,
		GzipEnabled:         &set.Nginx.GzipEnabled,
		TcpNoDelayEnabled:   &set.Nginx.TcpNoDelayEnabled,
		RuntimeUser:         &set.Nginx.RuntimeUser,
		Custom:              set.Nginx.Custom,
	}

	logRotationModel := &logRotationSettingsDto{
		Enabled:           &set.LogRotation.Enabled,
		MaximumLines:      &set.LogRotation.MaximumLines,
		IntervalUnit:      &set.LogRotation.IntervalUnit,
		IntervalUnitCount: &set.LogRotation.IntervalUnitCount,
	}

	certificateModel := &certificateAutoRenewSettingsDto{
		Enabled:           &set.CertificateAutoRenew.Enabled,
		IntervalUnit:      &set.CertificateAutoRenew.IntervalUnit,
		IntervalUnitCount: &set.CertificateAutoRenew.IntervalUnitCount,
	}

	bindingsModel := make([]bindingDto, 0)
	for _, b := range set.GlobalBindings {
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
