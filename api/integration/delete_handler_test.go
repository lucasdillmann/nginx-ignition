package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/integration"
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
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Delete(gomock.Any(), id).
				Return(nil)

			handler := deleteHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.DELETE("/api/integrations/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/integrations/"+id.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := deleteHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.DELETE("/api/integrations/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/integrations/invalid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Delete(gomock.Any(), id).
				Return(expectedErr)

			handler := deleteHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.DELETE("/api/integrations/:id", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/integrations/"+id.String(), nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
