package cache

import (
	"bytes"
	"encoding/json"
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

func Test_createHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 201 Created on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newCacheRequestDTO()
			commands := cache.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := createHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/caches", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/caches",
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusCreated, recorder.Code)
			var response cacheResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, payload.Name, response.Name)
			assert.NotEqual(t, uuid.Nil, response.ID)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			handler := createHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.POST("/api/caches", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/caches",
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

			payload := newCacheRequestDTO()
			expectedErr := assert.AnError
			commands := cache.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := createHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/caches", func(ginContext *gin.Context) {
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
				"POST",
				"/api/caches",
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
