package settings

import (
	"bytes"
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

func Test_putHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newSettingsDTO()
			commands := settings.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := putHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/settings", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PUT", "/api/settings", bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			handler := putHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/settings", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/settings",
				bytes.NewBufferString("invalid json"),
			)
			request.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newSettingsDTO()
			expectedErr := errors.New("update error")
			commands := settings.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := putHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/settings", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PUT", "/api/settings", bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
