package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func Test_PutHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := settingsDTO{
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
	body, _ := json.Marshal(payload)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := settings.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := putHandler{commands}
		r := gin.New()
		r.PUT("/api/settings", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/settings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("panics on invalid JSON", func(t *testing.T) {
		handler := putHandler{nil}
		r := gin.New()
		r.PUT("/api/settings", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/settings", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("update error")
		commands := settings.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		handler := putHandler{commands}
		r := gin.New()
		r.PUT("/api/settings", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/settings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
