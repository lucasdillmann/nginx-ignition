package accesslist

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
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
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Delete(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, idToDelete uuid.UUID) error {
					assert.Equal(t, id, idToDelete)
					return nil
				})

			engine := gin.New()
			handler := deleteHandler{
				commands: commands,
			}
			engine.DELETE("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/access-lists/"+id.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
			engine := gin.New()
			handler := deleteHandler{
				commands: nil,
			}
			engine.DELETE("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/access-lists/invalid-uuid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Delete(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			engine := gin.New()
			handler := deleteHandler{
				commands: commands,
			}
			engine.DELETE("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/access-lists/"+id.String(), nil)

			assert.PanicsWithValue(t, expectedErr, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
