package settings

import (
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func toDomain(
	nginx *nginxModel,
	logRotation *logRotationModel,
	certificate *certificateModel,
	bindings []bindingModel,
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
			TCPNoDelayEnabled:   nginx.TCPNoDelayEnabled,
			RuntimeUser:         nginx.RuntimeUser,
			Custom:              nginx.Custom,
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

func toBindingDomain(bindings []bindingModel) []binding.Binding {
	result := make([]binding.Binding, 0)

	for _, b := range bindings {
		result = append(result, binding.Binding{
			ID:            b.ID,
			Type:          binding.Type(b.Type),
			IP:            b.IP,
			Port:          b.Port,
			CertificateID: b.CertificateID,
		})
	}

	return result
}

func toModel(set *settings.Settings) (
	*nginxModel,
	*logRotationModel,
	*certificateModel,
	[]bindingModel,
	*buffersModel,
) {
	nginx := &nginxModel{
		ServerLogsEnabled:   set.Nginx.Logs.ServerLogsEnabled,
		ServerLogsLevel:     string(set.Nginx.Logs.ServerLogsLevel),
		AccessLogsEnabled:   set.Nginx.Logs.AccessLogsEnabled,
		ErrorLogsEnabled:    set.Nginx.Logs.ErrorLogsEnabled,
		ErrorLogsLevel:      string(set.Nginx.Logs.ErrorLogsLevel),
		ReadTimeout:         set.Nginx.Timeouts.Read,
		ConnectTimeout:      set.Nginx.Timeouts.Connect,
		SendTimeout:         set.Nginx.Timeouts.Send,
		KeepaliveTimeout:    set.Nginx.Timeouts.Keepalive,
		ClientBodyTimeout:   set.Nginx.Timeouts.ClientBody,
		WorkerProcesses:     set.Nginx.WorkerProcesses,
		WorkerConnections:   set.Nginx.WorkerConnections,
		DefaultContentType:  set.Nginx.DefaultContentType,
		ServerTokensEnabled: set.Nginx.ServerTokensEnabled,
		MaximumBodySizeMb:   set.Nginx.MaximumBodySizeMb,
		SendfileEnabled:     set.Nginx.SendfileEnabled,
		GzipEnabled:         set.Nginx.GzipEnabled,
		TCPNoDelayEnabled:   set.Nginx.TCPNoDelayEnabled,
		RuntimeUser:         set.Nginx.RuntimeUser,
		Custom:              set.Nginx.Custom,
	}

	logRotation := &logRotationModel{
		Enabled:           set.LogRotation.Enabled,
		MaximumLines:      set.LogRotation.MaximumLines,
		IntervalUnit:      string(set.LogRotation.IntervalUnit),
		IntervalUnitCount: set.LogRotation.IntervalUnitCount,
	}

	certificate := &certificateModel{
		Enabled:           set.CertificateAutoRenew.Enabled,
		IntervalUnit:      string(set.CertificateAutoRenew.IntervalUnit),
		IntervalUnitCount: set.CertificateAutoRenew.IntervalUnitCount,
	}

	bindings := toBindingModel(set.GlobalBindings)

	buffers := &buffersModel{
		ClientBodyKb:            set.Nginx.Buffers.ClientBodyKb,
		ClientHeaderKb:          set.Nginx.Buffers.ClientHeaderKb,
		LargeClientHeaderSizeKb: set.Nginx.Buffers.LargeClientHeader.SizeKb,
		LargeClientHeaderAmount: set.Nginx.Buffers.LargeClientHeader.Amount,
		OutputSizeKb:            set.Nginx.Buffers.Output.SizeKb,
		OutputAmount:            set.Nginx.Buffers.Output.Amount,
	}

	return nginx, logRotation, certificate, bindings, buffers
}

func toBindingModel(bindings []binding.Binding) []bindingModel {
	result := make([]bindingModel, 0)

	for _, b := range bindings {
		result = append(result, bindingModel{
			ID:            b.ID,
			Type:          string(b.Type),
			IP:            b.IP,
			Port:          b.Port,
			CertificateID: b.CertificateID,
		})
	}

	return result
}
