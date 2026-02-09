package nginx

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_service_GetTrafficStats(t *testing.T) {
	t.Run("returns stats when nginx returns valid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{
			Nginx: &settings.NginxSettings{
				Stats: &settings.NginxStatsSettings{
					Enabled: true,
				},
			},
		}, nil)

		mockJSON := `{"hostName": "test-host", "connections": {"active": 10}}`
		client := &http.Client{
			Transport: &mockTransport{
				RoundTripFunc: func(_ *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(mockJSON)),
						Header:     make(http.Header),
					}, nil
				},
			},
		}

		nginxService := &service{
			settingsCommands: settingsCmds,
			semaphore: &semaphore{
				state: runningState,
			},
			statsClient: client,
		}

		stats, err := nginxService.GetTrafficStats(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, "test-host", stats.HostName)
		assert.Equal(t, uint64(10), stats.Connections.Active)
	})

	t.Run("returns error when stats not enabled in settings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{
			Nginx: &settings.NginxSettings{
				Stats: &settings.NginxStatsSettings{
					Enabled: false,
				},
			},
		}, nil)

		nginxService := &service{
			settingsCommands: settingsCmds,
		}

		stats, err := nginxService.GetTrafficStats(context.Background())

		assert.Nil(t, stats)
		assert.Error(t, err)
		var coreErr *coreerror.CoreError
		assert.ErrorAs(t, err, &coreErr)
		assert.Equal(t, i18n.K.CoreNginxStatsNotEnabled, coreErr.Message.Key)
	})

	t.Run("returns error when stats settings is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{
			Nginx: &settings.NginxSettings{
				Stats: nil,
			},
		}, nil)

		nginxService := &service{
			settingsCommands: settingsCmds,
		}

		stats, err := nginxService.GetTrafficStats(context.Background())

		assert.Nil(t, stats)
		assert.Error(t, err)
		var coreErr *coreerror.CoreError
		assert.ErrorAs(t, err, &coreErr)
		assert.Equal(t, i18n.K.CoreNginxStatsNotEnabled, coreErr.Message.Key)
	})

	t.Run("returns error when nginx settings is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{
			Nginx: nil,
		}, nil)

		nginxService := &service{
			settingsCommands: settingsCmds,
		}

		stats, err := nginxService.GetTrafficStats(context.Background())

		assert.Nil(t, stats)
		assert.Error(t, err)
		var coreErr *coreerror.CoreError
		assert.ErrorAs(t, err, &coreErr)
		assert.Equal(t, i18n.K.CoreNginxStatsNotEnabled, coreErr.Message.Key)
	})

	t.Run("returns error when settings command fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(gomock.Any()).Return(nil, expectedErr)

		nginxService := &service{
			settingsCommands: settingsCmds,
		}

		stats, err := nginxService.GetTrafficStats(context.Background())

		assert.Nil(t, stats)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("returns error when nginx is not running", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		settingsCmds := settings.NewMockedCommands(ctrl)
		settingsCmds.EXPECT().Get(gomock.Any()).Return(&settings.Settings{
			Nginx: &settings.NginxSettings{
				Stats: &settings.NginxStatsSettings{
					Enabled: true,
				},
			},
		}, nil)

		nginxService := &service{
			settingsCommands: settingsCmds,
			semaphore: &semaphore{
				state: stoppedState,
			},
		}

		stats, err := nginxService.GetTrafficStats(context.Background())

		assert.Nil(t, stats)
		assert.Error(t, err)
		var coreErr *coreerror.CoreError
		assert.ErrorAs(t, err, &coreErr)
		assert.Equal(t, i18n.K.CoreNginxNotRunning, coreErr.Message.Key)
	})
}
