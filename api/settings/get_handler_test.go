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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with settings data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			subject := newSettings()
			commands := settings.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any()).
				Return(subject, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/settings", nil)

			handler := getHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response settingsDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, subject.Nginx.GzipEnabled, *response.Nginx.GzipEnabled)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := errors.New("settings error")
			commands := settings.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest("GET", "/api/settings", nil)

			handler := getHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
