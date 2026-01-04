package host

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/host"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_updateHandler(t *testing.T) {
	id := uuid.New()

	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newHostRequestDTO()
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			body, _ := json.Marshal(payload)
			request := httptest.NewRequest("PUT", "/api/hosts/"+id.String(), bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			payload := newHostRequestDTO()
			body, _ := json.Marshal(payload)
			request := httptest.NewRequest("PUT", "/api/hosts/invalid", bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/hosts/"+id.String(),
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

			payload := newHostRequestDTO()
			expectedErr := errors.New("update error")
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/hosts/:id", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			body, _ := json.Marshal(payload)
			request := httptest.NewRequest("PUT", "/api/hosts/"+id.String(), bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
