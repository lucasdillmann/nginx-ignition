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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_logsHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with logs on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			logs := []string{"log line 1", "log line 2"}
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetHostLogs(gomock.Any(), id, "access", 50).
				Return(logs, nil)

			handler := logsHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/hosts/"+id.String()+"/logs/access", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []string
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, logs, response)
		})

		t.Run("returns 400 Bad Request on invalid line count", func(t *testing.T) {
			handler := logsHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/hosts/"+uuid.New().String()+"/logs/access?lines=abc",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid qualifier", func(t *testing.T) {
			handler := logsHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/hosts/"+uuid.New().String()+"/logs/invalid",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := logsHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.GET("/api/hosts/:id/logs/:qualifier", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/hosts/invalid/logs/access", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
