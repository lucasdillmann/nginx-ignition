package cache

import (
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

func Test_deleteHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := cache.NewMockedCommands(controller)
			commands.EXPECT().
				Delete(gomock.Any(), id).
				Return(nil)

			handler := deleteHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.DELETE("/api/caches/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/caches/"+id.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := deleteHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.DELETE("/api/caches/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/caches/invalid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			expectedErr := errors.New("delete error")
			commands := cache.NewMockedCommands(controller)
			commands.EXPECT().
				Delete(gomock.Any(), id).
				Return(expectedErr)

			handler := deleteHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.DELETE("/api/caches/:id", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/caches/"+id.String(), nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
