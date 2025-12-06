package settings

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func toDomain(
	nginx *nginxModel,
	logRotation *logRotationModel,
	certificate *certificateModel,
	bindings []*bindingModel,
	buffers *buffersModel,
) *settings.Settings {
	return &settings.Settings{
		Nginx: &settings.NginxSettings{
			Logs: &settings.NginxLogsSettings{
				ServerLogsEnabled: nginx.ServerLogsEnabled,
				ServerLogsLevel:   settings.LogLevel(nginx.ServerLogsLevel),
				AccessLogsEnabled: nginx.AccessLogsEnabled,
				ErrorLogsEnabled:  nginx.ErrorLogsEnabled,
				ErrorLogsLevel:    settings.LogLevel(nginx.ErrorLogsLevel),
			},
			Timeouts: &settings.NginxTimeoutsSettings{
				Read:       nginx.ReadTimeout,
				Connect:    nginx.ConnectTimeout,
				Send:       nginx.SendTimeout,
				Keepalive:  nginx.KeepaliveTimeout,
				ClientBody: nginx.ClientBodyTimeout,
			},
			Buffers: &settings.NginxBuffersSettings{
				ClientBodyKb:   buffers.ClientBodyKb,
				ClientHeaderKb: buffers.ClientHeaderKb,
				LargeClientHeader: &settings.NginxBufferSize{
					SizeKb: buffers.LargeClientHeaderSizeKb,
					Amount: buffers.LargeClientHeaderAmount,
				},
				Output: &settings.NginxBufferSize{
					SizeKb: buffers.OutputSizeKb,
					Amount: buffers.OutputAmount,
				},
			},
			WorkerProcesses:     nginx.WorkerProcesses,
			WorkerConnections:   nginx.WorkerConnections,
			DefaultContentType:  nginx.DefaultContentType,
			ServerTokensEnabled: nginx.ServerTokensEnabled,
			MaximumBodySizeMb:   nginx.MaximumBodySizeMb,
			SendfileEnabled:     nginx.SendfileEnabled,
			GzipEnabled:         nginx.GzipEnabled,
			TcpNoDelayEnabled:   nginx.TcpNoDelayEnabled,
			RuntimeUser:         settings.RuntimeUser(nginx.RuntimeUser),
		},
		LogRotation: &settings.LogRotationSettings{
			Enabled:           logRotation.Enabled,
			MaximumLines:      logRotation.MaximumLines,
			IntervalUnit:      settings.TimeUnit(logRotation.IntervalUnit),
			IntervalUnitCount: logRotation.IntervalUnitCount,
		},
		CertificateAutoRenew: &settings.CertificateAutoRenewSettings{
			Enabled:           certificate.Enabled,
			IntervalUnit:      settings.TimeUnit(certificate.IntervalUnit),
			IntervalUnitCount: certificate.IntervalUnitCount,
		},
		GlobalBindings: toBindingDomain(bindings),
	}
}

func toBindingDomain(bindings []*bindingModel) []*host.Binding {
	var result []*host.Binding

	for _, binding := range bindings {
		result = append(result, &host.Binding{
			ID:            binding.ID,
			Type:          host.BindingType(binding.Type),
			IP:            binding.IP,
			Port:          binding.Port,
			CertificateID: binding.CertificateID,
		})
	}

	return result
}

func toModel(settings *settings.Settings) (
	*nginxModel,
	*logRotationModel,
	*certificateModel,
	[]*bindingModel,
	*buffersModel,
) {
	nginx := &nginxModel{
		ServerLogsEnabled:   settings.Nginx.Logs.ServerLogsEnabled,
		ServerLogsLevel:     string(settings.Nginx.Logs.ServerLogsLevel),
		AccessLogsEnabled:   settings.Nginx.Logs.AccessLogsEnabled,
		ErrorLogsEnabled:    settings.Nginx.Logs.ErrorLogsEnabled,
		ErrorLogsLevel:      string(settings.Nginx.Logs.ErrorLogsLevel),
		ReadTimeout:         settings.Nginx.Timeouts.Read,
		ConnectTimeout:      settings.Nginx.Timeouts.Connect,
		SendTimeout:         settings.Nginx.Timeouts.Send,
		KeepaliveTimeout:    settings.Nginx.Timeouts.Keepalive,
		ClientBodyTimeout:   settings.Nginx.Timeouts.ClientBody,
		WorkerProcesses:     settings.Nginx.WorkerProcesses,
		WorkerConnections:   settings.Nginx.WorkerConnections,
		DefaultContentType:  settings.Nginx.DefaultContentType,
		ServerTokensEnabled: settings.Nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   settings.Nginx.MaximumBodySizeMb,
		SendfileEnabled:     settings.Nginx.SendfileEnabled,
		GzipEnabled:         settings.Nginx.GzipEnabled,
		TcpNoDelayEnabled:   settings.Nginx.TcpNoDelayEnabled,
		RuntimeUser:         string(settings.Nginx.RuntimeUser),
	}

	logRotation := &logRotationModel{
		Enabled:           settings.LogRotation.Enabled,
		MaximumLines:      settings.LogRotation.MaximumLines,
		IntervalUnit:      string(settings.LogRotation.IntervalUnit),
		IntervalUnitCount: settings.LogRotation.IntervalUnitCount,
	}

	certificate := &certificateModel{
		Enabled:           settings.CertificateAutoRenew.Enabled,
		IntervalUnit:      string(settings.CertificateAutoRenew.IntervalUnit),
		IntervalUnitCount: settings.CertificateAutoRenew.IntervalUnitCount,
	}

	bindings := toBindingModel(settings.GlobalBindings)

	buffers := &buffersModel{
		ClientBodyKb:            settings.Nginx.Buffers.ClientBodyKb,
		ClientHeaderKb:          settings.Nginx.Buffers.ClientHeaderKb,
		LargeClientHeaderSizeKb: settings.Nginx.Buffers.LargeClientHeader.SizeKb,
		LargeClientHeaderAmount: settings.Nginx.Buffers.LargeClientHeader.Amount,
		OutputSizeKb:            settings.Nginx.Buffers.Output.SizeKb,
		OutputAmount:            settings.Nginx.Buffers.Output.Amount,
	}

	return nginx, logRotation, certificate, bindings, buffers
}

func toBindingModel(bindings []*host.Binding) []*bindingModel {
	var result []*bindingModel

	for _, binding := range bindings {
		result = append(result, &bindingModel{
			ID:            binding.ID,
			Type:          string(binding.Type), // Assuming host.BindingType has a String method
			IP:            binding.IP,
			Port:          binding.Port,
			CertificateID: binding.CertificateID,
		})
	}

	return result
}
