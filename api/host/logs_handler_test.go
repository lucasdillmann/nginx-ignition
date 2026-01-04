package host

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

func Test_LogsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with logs on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			mockLogs := []string{"log line 1", "log line 2"}
			commands := nginx.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetHostLogs(gomock.Any(), id, "access", 50).
				Return(mockLogs, nil)

			handler := logsHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/hosts/"+id.String()+"/logs/access", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp []string
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, mockLogs, resp)
		})

		t.Run("returns 400 Bad Request on invalid line count", func(t *testing.T) {
			handler := logsHandler{
				commands: nil,
			}
			r := gin.New()
			r.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/api/hosts/"+uuid.New().String()+"/logs/access?lines=abc",
				nil,
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("returns 404 Not Found on invalid qualifier", func(t *testing.T) {
			handler := logsHandler{
				commands: nil,
			}
			r := gin.New()
			r.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/api/hosts/"+uuid.New().String()+"/logs/invalid",
				nil,
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := logsHandler{
				commands: nil,
			}
			r := gin.New()
			r.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/hosts/invalid/logs/access", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	})
}
