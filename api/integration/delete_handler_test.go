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

func Test_DeleteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				Delete(gomock.Any(), id).
				Return(nil)

			handler := deleteHandler{
				commands: commands,
			}
			r := gin.New()
			r.DELETE("/api/integrations/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/integrations/"+id.String(), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := deleteHandler{
				commands: nil,
			}
			r := gin.New()
			r.DELETE("/api/integrations/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/integrations/invalid", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				Delete(gomock.Any(), id).
				Return(expectedErr)

			handler := deleteHandler{
				commands: commands,
			}
			r := gin.New()
			r.DELETE("/api/integrations/:id", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/integrations/"+id.String(), nil)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
