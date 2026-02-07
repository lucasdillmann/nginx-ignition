package settings

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func Test_toDTO(t *testing.T) {
	t.Run("returns nil when input is nil", func(t *testing.T) {
		result := toDTO(nil)
		assert.Nil(t, result)
	})

	t.Run("converts domain object to DTO", func(t *testing.T) {
		subject := newSettings()
		subject.Nginx.Custom = ptr.Of("custom-config")
		subject.Nginx.Stats.DatabaseLocation = ptr.Of("/var/lib/nginx/stats.db")
		subject.GlobalBindings = []binding.Binding{
			{
				Type:          binding.HTTPBindingType,
				IP:            "127.0.0.1",
				Port:          8080,
				CertificateID: ptr.Of(uuid.New()),
			},
		}

		result := toDTO(subject)
		assert.NotNil(t, result)

		// Nginx
		nginxSubject := subject.Nginx
		nginxResult := result.Nginx

		assert.Equal(t, nginxSubject.WorkerProcesses, *nginxResult.WorkerProcesses)
		assert.Equal(t, nginxSubject.WorkerConnections, *nginxResult.WorkerConnections)
		assert.Equal(t, nginxSubject.DefaultContentType, *nginxResult.DefaultContentType)
		assert.Equal(t, nginxSubject.ServerTokensEnabled, *nginxResult.ServerTokensEnabled)
		assert.Equal(t, nginxSubject.MaximumBodySizeMb, *nginxResult.MaximumBodySizeMb)
		assert.Equal(t, nginxSubject.SendfileEnabled, *nginxResult.SendfileEnabled)
		assert.Equal(t, nginxSubject.GzipEnabled, *nginxResult.GzipEnabled)
		assert.Equal(t, nginxSubject.TCPNoDelayEnabled, *nginxResult.TCPNoDelayEnabled)
		assert.Equal(t, nginxSubject.RuntimeUser, *nginxResult.RuntimeUser)
		assert.Equal(t, nginxSubject.Custom, nginxResult.Custom)

		// Nginx Logs
		logsSubject := nginxSubject.Logs
		logsResult := nginxResult.Logs

		assert.Equal(t, logsSubject.ServerLogsEnabled, *logsResult.ServerLogsEnabled)
		assert.Equal(t, logsSubject.ServerLogsLevel, *logsResult.ServerLogsLevel)
		assert.Equal(t, logsSubject.AccessLogsEnabled, *logsResult.AccessLogsEnabled)
		assert.Equal(t, logsSubject.ErrorLogsEnabled, *logsResult.ErrorLogsEnabled)
		assert.Equal(t, logsSubject.ErrorLogsLevel, *logsResult.ErrorLogsLevel)

		// Nginx Timeouts
		timeoutsSubject := nginxSubject.Timeouts
		timeoutsResult := nginxResult.Timeouts

		assert.Equal(t, timeoutsSubject.Read, *timeoutsResult.Read)
		assert.Equal(t, timeoutsSubject.Connect, *timeoutsResult.Connect)
		assert.Equal(t, timeoutsSubject.Send, *timeoutsResult.Send)
		assert.Equal(t, timeoutsSubject.Keepalive, *timeoutsResult.Keepalive)
		assert.Equal(t, timeoutsSubject.ClientBody, *timeoutsResult.ClientBody)

		// Nginx Buffers
		buffersSubject := nginxSubject.Buffers
		buffersResult := nginxResult.Buffers

		assert.Equal(t, buffersSubject.ClientBodyKb, *buffersResult.ClientBodyKb)
		assert.Equal(t, buffersSubject.ClientHeaderKb, *buffersResult.ClientHeaderKb)
		assert.Equal(
			t,
			buffersSubject.LargeClientHeader.SizeKb,
			*buffersResult.LargeClientHeader.SizeKb,
		)
		assert.Equal(
			t,
			buffersSubject.LargeClientHeader.Amount,
			*buffersResult.LargeClientHeader.Amount,
		)
		assert.Equal(t, buffersSubject.Output.SizeKb, *buffersResult.Output.SizeKb)
		assert.Equal(t, buffersSubject.Output.Amount, *buffersResult.Output.Amount)

		// Nginx Stats
		statsSubject := nginxSubject.Stats
		statsResult := nginxResult.Stats

		assert.Equal(t, statsSubject.Enabled, *statsResult.Enabled)
		assert.Equal(t, statsSubject.Persistent, *statsResult.Persistent)
		assert.Equal(t, statsSubject.MaximumSizeMB, *statsResult.MaximumSizeMB)
		assert.Equal(t, statsSubject.DatabaseLocation, statsResult.DatabaseLocation)

		// Log Rotation
		logRotationSubject := subject.LogRotation
		logRotationResult := result.LogRotation

		assert.Equal(t, logRotationSubject.Enabled, *logRotationResult.Enabled)
		assert.Equal(t, logRotationSubject.MaximumLines, *logRotationResult.MaximumLines)
		assert.Equal(t, logRotationSubject.IntervalUnit, *logRotationResult.IntervalUnit)
		assert.Equal(t, logRotationSubject.IntervalUnitCount, *logRotationResult.IntervalUnitCount)

		// Certificate Auto Renew
		certSubject := subject.CertificateAutoRenew
		certResult := result.CertificateAutoRenew

		assert.Equal(t, certSubject.Enabled, *certResult.Enabled)
		assert.Equal(t, certSubject.IntervalUnit, *certResult.IntervalUnit)
		assert.Equal(t, certSubject.IntervalUnitCount, *certResult.IntervalUnitCount)

		// Global Bindings
		assert.Len(t, result.GlobalBindings, 1)
		bindingSubject := subject.GlobalBindings[0]
		bindingResult := result.GlobalBindings[0]

		assert.Equal(t, bindingSubject.Type, *bindingResult.Type)
		assert.Equal(t, bindingSubject.IP, *bindingResult.IP)
		assert.Equal(t, bindingSubject.Port, *bindingResult.Port)
		assert.Equal(t, bindingSubject.CertificateID, bindingResult.CertificateID)
	})
}

func Test_toDomain(t *testing.T) {
	t.Run("returns nil when critical fields are missing", func(t *testing.T) {
		assert.Nil(t, toDomain(&settingsDTO{}))
		assert.Nil(t, toDomain(&settingsDTO{Nginx: &nginxSettingsDTO{}}))
		assert.Nil(t, toDomain(&settingsDTO{LogRotation: &logRotationSettingsDTO{}}))
		assert.Nil(
			t,
			toDomain(&settingsDTO{CertificateAutoRenew: &certificateAutoRenewSettingsDTO{}}),
		)
	})

	t.Run("converts DTO to domain object", func(t *testing.T) {
		payload := newSettingsDTO()
		payload.Nginx.Custom = ptr.Of("custom-config")
		payload.Nginx.Stats.DatabaseLocation = ptr.Of("/var/lib/nginx/stats.db")
		payload.GlobalBindings = []bindingDTO{
			{
				Type:          ptr.Of(binding.HTTPSBindingType),
				IP:            ptr.Of("0.0.0.0"),
				Port:          ptr.Of(443),
				CertificateID: ptr.Of(uuid.New()),
			},
		}

		result := toDomain(payload)
		assert.NotNil(t, result)

		// Nginx
		nginxPayload := payload.Nginx
		nginxResult := result.Nginx

		assert.Equal(t, *nginxPayload.WorkerProcesses, nginxResult.WorkerProcesses)
		assert.Equal(t, *nginxPayload.WorkerConnections, nginxResult.WorkerConnections)
		assert.Equal(t, *nginxPayload.DefaultContentType, nginxResult.DefaultContentType)
		assert.Equal(t, *nginxPayload.ServerTokensEnabled, nginxResult.ServerTokensEnabled)
		assert.Equal(t, *nginxPayload.MaximumBodySizeMb, nginxResult.MaximumBodySizeMb)
		assert.Equal(t, *nginxPayload.SendfileEnabled, nginxResult.SendfileEnabled)
		assert.Equal(t, *nginxPayload.GzipEnabled, nginxResult.GzipEnabled)
		assert.Equal(t, *nginxPayload.TCPNoDelayEnabled, nginxResult.TCPNoDelayEnabled)
		assert.Equal(t, *nginxPayload.RuntimeUser, nginxResult.RuntimeUser)
		assert.Equal(t, nginxPayload.Custom, nginxResult.Custom)

		// Nginx Logs
		logsPayload := nginxPayload.Logs
		logsResult := nginxResult.Logs

		assert.Equal(t, *logsPayload.ServerLogsEnabled, logsResult.ServerLogsEnabled)
		assert.Equal(t, *logsPayload.ServerLogsLevel, logsResult.ServerLogsLevel)
		assert.Equal(t, *logsPayload.AccessLogsEnabled, logsResult.AccessLogsEnabled)
		assert.Equal(t, *logsPayload.ErrorLogsEnabled, logsResult.ErrorLogsEnabled)
		assert.Equal(t, *logsPayload.ErrorLogsLevel, logsResult.ErrorLogsLevel)

		// Nginx Timeouts
		timeoutsPayload := nginxPayload.Timeouts
		timeoutsResult := nginxResult.Timeouts

		assert.Equal(t, *timeoutsPayload.Read, timeoutsResult.Read)
		assert.Equal(t, *timeoutsPayload.Connect, timeoutsResult.Connect)
		assert.Equal(t, *timeoutsPayload.Send, timeoutsResult.Send)
		assert.Equal(t, *timeoutsPayload.Keepalive, timeoutsResult.Keepalive)
		assert.Equal(t, *timeoutsPayload.ClientBody, timeoutsResult.ClientBody)

		// Nginx Buffers
		buffersPayload := nginxPayload.Buffers
		buffersResult := nginxResult.Buffers

		assert.Equal(t, *buffersPayload.ClientBodyKb, buffersResult.ClientBodyKb)
		assert.Equal(t, *buffersPayload.ClientHeaderKb, buffersResult.ClientHeaderKb)
		assert.Equal(
			t,
			*buffersPayload.LargeClientHeader.SizeKb,
			buffersResult.LargeClientHeader.SizeKb,
		)
		assert.Equal(
			t,
			*buffersPayload.LargeClientHeader.Amount,
			buffersResult.LargeClientHeader.Amount,
		)
		assert.Equal(t, *buffersPayload.Output.SizeKb, buffersResult.Output.SizeKb)
		assert.Equal(t, *buffersPayload.Output.Amount, buffersResult.Output.Amount)

		// Nginx Stats
		statsPayload := nginxPayload.Stats
		statsResult := nginxResult.Stats

		assert.Equal(t, *statsPayload.Enabled, statsResult.Enabled)
		assert.Equal(t, *statsPayload.Persistent, statsResult.Persistent)
		assert.Equal(t, *statsPayload.MaximumSizeMB, statsResult.MaximumSizeMB)
		assert.Equal(t, statsPayload.DatabaseLocation, statsResult.DatabaseLocation)

		// Log Rotation
		logRotationPayload := payload.LogRotation
		logRotationResult := result.LogRotation

		assert.Equal(t, *logRotationPayload.Enabled, logRotationResult.Enabled)
		assert.Equal(t, *logRotationPayload.MaximumLines, logRotationResult.MaximumLines)
		assert.Equal(t, *logRotationPayload.IntervalUnit, logRotationResult.IntervalUnit)
		assert.Equal(t, *logRotationPayload.IntervalUnitCount, logRotationResult.IntervalUnitCount)

		// Certificate Auto Renew
		certPayload := payload.CertificateAutoRenew
		certResult := result.CertificateAutoRenew

		assert.Equal(t, *certPayload.Enabled, certResult.Enabled)
		assert.Equal(t, *certPayload.IntervalUnit, certResult.IntervalUnit)
		assert.Equal(t, *certPayload.IntervalUnitCount, certResult.IntervalUnitCount)

		// Global Bindings
		assert.Len(t, result.GlobalBindings, 1)
		bindingPayload := payload.GlobalBindings[0]
		bindingResult := result.GlobalBindings[0]

		assert.NotEqual(t, uuid.Nil, bindingResult.ID)
		assert.Equal(t, *bindingPayload.Type, bindingResult.Type)
		assert.Equal(t, *bindingPayload.IP, bindingResult.IP)
		assert.Equal(t, *bindingPayload.Port, bindingResult.Port)
		assert.Equal(t, bindingPayload.CertificateID, bindingResult.CertificateID)
	})

	t.Run("converts empty database location to nil", func(t *testing.T) {
		payload := newSettingsDTO()
		payload.Nginx.Stats.DatabaseLocation = ptr.Of("   ")
		result := toDomain(payload)

		assert.NotNil(t, result)
		assert.Nil(t, result.Nginx.Stats.DatabaseLocation)
	})
}
