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
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_reloadHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				Reload(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := reloadHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/nginx/reload", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/nginx/reload", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 424 Failed Dependency on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				Reload(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := reloadHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/nginx/reload", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/nginx/reload", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusFailedDependency, recorder.Code)
			var response map[string]string
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, expectedErr.Error(), response["message"])
		})
	})
}
