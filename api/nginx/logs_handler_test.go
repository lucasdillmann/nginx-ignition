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

func Test_logsHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with logs on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			logs := []string{"log line 1", "log line 2"}
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetMainLogs(gomock.Any(), 50).
				Return(logs, nil)

			handler := logsHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/logs", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/logs", nil)
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
			engine.GET("/api/nginx/logs", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/logs?lines=abc", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetMainLogs(gomock.Any(), 50).
				Return(nil, expectedErr)

			handler := logsHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/logs", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/logs", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
