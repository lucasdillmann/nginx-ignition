package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_Converter(t *testing.T) {
	t.Run("toDTO", func(t *testing.T) {
		t.Run("returns nil when input is nil", func(t *testing.T) {
			result := toDTO(nil)
			assert.Nil(t, result)
		})

		t.Run("converts domain object to DTO", func(t *testing.T) {
			input := &settings.Settings{
				Nginx: &settings.NginxSettings{
					GzipEnabled: true,
					Logs:        &settings.NginxLogsSettings{},
					Timeouts:    &settings.NginxTimeoutsSettings{},
					Buffers: &settings.NginxBuffersSettings{
						LargeClientHeader: &settings.NginxBufferSize{},
						Output:            &settings.NginxBufferSize{},
					},
				},
				LogRotation:          &settings.LogRotationSettings{},
				CertificateAutoRenew: &settings.CertificateAutoRenewSettings{},
			}

			result := toDTO(input)

			assert.NotNil(t, result)
			assert.Equal(t, input.Nginx.GzipEnabled, *result.Nginx.GzipEnabled)
		})
	})

	t.Run("toDomain", func(t *testing.T) {
		t.Run("returns nil when critical fields are missing", func(t *testing.T) {
			result := toDomain(&settingsDTO{})
			assert.Nil(t, result)
		})

		t.Run("converts DTO to domain object", func(t *testing.T) {
			input := &settingsDTO{
				Nginx: &nginxSettingsDTO{
					GzipEnabled:         ptr.Of(true),
					ServerTokensEnabled: ptr.Of(true),
					SendfileEnabled:     ptr.Of(true),
					TCPNoDelayEnabled:   ptr.Of(true),
					WorkerProcesses:     ptr.Of(0),
					WorkerConnections:   ptr.Of(0),
					MaximumBodySizeMb:   ptr.Of(0),
					DefaultContentType:  ptr.Of(""),
					RuntimeUser:         ptr.Of(""),
					Logs: &nginxLogsSettingsDTO{
						ServerLogsEnabled: ptr.Of(true),
						AccessLogsEnabled: ptr.Of(true),
						ErrorLogsEnabled:  ptr.Of(true),
						ServerLogsLevel:   ptr.Of(settings.WarnLogLevel),
						ErrorLogsLevel:    ptr.Of(settings.WarnLogLevel),
					},
					Timeouts: &nginxTimeoutsSettingsDTO{
						Read:       ptr.Of(0),
						Connect:    ptr.Of(0),
						Send:       ptr.Of(0),
						Keepalive:  ptr.Of(0),
						ClientBody: ptr.Of(0),
					},
					Buffers: &nginxBuffersSettingsDTO{
						ClientBodyKb:   ptr.Of(0),
						ClientHeaderKb: ptr.Of(0),
						LargeClientHeader: &nginxBufferSizeDTO{
							SizeKb: ptr.Of(0),
							Amount: ptr.Of(0),
						},
						Output: &nginxBufferSizeDTO{
							SizeKb: ptr.Of(0),
							Amount: ptr.Of(0),
						},
					},
				},
				LogRotation: &logRotationSettingsDTO{
					Enabled:           ptr.Of(true),
					MaximumLines:      ptr.Of(0),
					IntervalUnit:      ptr.Of(settings.MinutesTimeUnit),
					IntervalUnitCount: ptr.Of(0),
				},
				CertificateAutoRenew: &certificateAutoRenewSettingsDTO{
					Enabled:           ptr.Of(true),
					IntervalUnit:      ptr.Of(settings.MinutesTimeUnit),
					IntervalUnitCount: ptr.Of(0),
				},
			}

			result := toDomain(input)

			assert.NotNil(t, result)
			assert.Equal(t, *input.Nginx.GzipEnabled, result.Nginx.GzipEnabled)
		})
	})
}
