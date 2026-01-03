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

func Test_LogsHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with logs on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockLogs := []string{"log line 1", "log line 2"}
		commands := nginx.NewMockedCommands(ctrl)
		commands.EXPECT().
			GetMainLogs(gomock.Any(), 50).
			Return(mockLogs, nil)

		handler := logsHandler{commands}
		r := gin.New()
		r.GET("/api/nginx/logs", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/nginx/logs", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, mockLogs, resp)
	})

	t.Run("returns 400 Bad Request on invalid line count", func(t *testing.T) {
		handler := logsHandler{nil}
		r := gin.New()
		r.GET("/api/nginx/logs", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/nginx/logs?lines=abc", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := nginx.NewMockedCommands(ctrl)
		commands.EXPECT().
			GetMainLogs(gomock.Any(), 50).
			Return(nil, expectedErr)

		handler := logsHandler{commands}
		r := gin.New()
		r.GET("/api/nginx/logs", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/nginx/logs", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
