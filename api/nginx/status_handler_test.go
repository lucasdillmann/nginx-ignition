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

func Test_statusHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with running status", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetStatus(gomock.Any()).
				Return(true)

			handler := statusHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/status", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/status", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response map[string]bool
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.True(t, response["running"])
		})
	})
}
