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

func Test_StopHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := nginx.NewMockedCommands(ctrl)
			commands.EXPECT().
				Stop(gomock.Any()).
				Return(nil)

			handler := stopHandler{
				commands: commands,
			}
			r := gin.New()
			r.POST("/api/nginx/stop", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/nginx/stop", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})

		t.Run("returns 424 Failed Dependency on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(ctrl)
			commands.EXPECT().
				Stop(gomock.Any()).
				Return(expectedErr)

			handler := stopHandler{
				commands: commands,
			}
			r := gin.New()
			r.POST("/api/nginx/stop", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/nginx/stop", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusFailedDependency, w.Code)
			var resp map[string]string
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, expectedErr.Error(), resp["message"])
		})
	})
}
