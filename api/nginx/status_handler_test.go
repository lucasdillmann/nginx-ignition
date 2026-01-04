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

func Test_StatusHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with running status", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := nginx.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetStatus(gomock.Any()).
				Return(true)

			handler := statusHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/nginx/status", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/nginx/status", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp map[string]bool
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.True(t, resp["running"])
		})
	})
}
