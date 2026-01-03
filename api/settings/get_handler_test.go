package settings

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_GetHandler_Handle(t *testing.T) {
	t.Run("returns 200 OK with settings data on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSettings := &settings.Settings{
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

		commands := settings.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any()).
			Return(mockSettings, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/settings", nil)

		handler := getHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp settingsDTO
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, mockSettings.Nginx.GzipEnabled, *resp.Nginx.GzipEnabled)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("settings error")
		commands := settings.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any()).
			Return(nil, expectedErr)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/settings", nil)

		handler := getHandler{commands}
		assert.PanicsWithValue(t, expectedErr, func() {
			handler.handle(ctx)
		})
	})
}
