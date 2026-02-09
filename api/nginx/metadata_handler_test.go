package nginx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/nginx"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_metadataHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with metadata", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			metadataData := newMetadata()
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetMetadata(gomock.Any()).
				Return(metadataData, nil)

			settingsData := newSettings()
			settingsCommands := settings.NewMockedCommands(controller)
			settingsCommands.EXPECT().
				Get(gomock.Any()).
				Return(settingsData, nil)

			handler := metadataHandler{
				nginxCommands:    commands,
				settingsCommands: settingsCommands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/metadata", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/metadata", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response map[string]any
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, metadataData.Version, response["version"])

			statsResponse := response["stats"].(map[string]any)
			assert.Equal(t, settingsData.Nginx.Stats.Enabled, statsResponse["enabled"])
			assert.Equal(t, settingsData.Nginx.Stats.AllHosts, statsResponse["allHosts"])
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetMetadata(gomock.Any()).
				Return(nil, expectedErr)

			handler := metadataHandler{
				nginxCommands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/metadata", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/metadata", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
