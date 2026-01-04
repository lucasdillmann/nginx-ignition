package cache

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

	"dillmann.com.br/nginx-ignition/core/cache"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_updateHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := newCacheRequestDTO()
			commands := cache.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/caches/:id", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/caches/"+id.String(),
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			payload := newCacheRequestDTO()
			handler := updateHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/caches/:id", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/caches/invalid",
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			id := uuid.New()
			handler := updateHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/caches/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/caches/"+id.String(),
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

			id := uuid.New()
			payload := newCacheRequestDTO()
			expectedErr := errors.New("update error")
			commands := cache.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/caches/:id", func(ginContext *gin.Context) {
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
			request := httptest.NewRequest(
				"PUT",
				"/api/caches/"+id.String(),
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
