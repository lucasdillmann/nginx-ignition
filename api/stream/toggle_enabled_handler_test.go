package stream

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_ToggleEnabledHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			mockStream := &stream.Stream{
				ID:      id,
				Enabled: true,
			}
			commands := stream.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(mockStream, nil)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			r := gin.New()
			r.PATCH("/api/streams/:id/toggle-enabled", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/streams/"+id.String()+"/toggle-enabled", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
			assert.False(t, mockStream.Enabled)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := toggleEnabledHandler{
				commands: nil,
			}
			r := gin.New()
			r.PATCH("/api/streams/:id/toggle-enabled", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/streams/invalid/toggle-enabled", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("returns 404 Not Found when stream does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			commands := stream.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, nil)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			r := gin.New()
			r.PATCH("/api/streams/:id/toggle-enabled", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/streams/"+id.String()+"/toggle-enabled", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := stream.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			r := gin.New()
			r.PATCH("/api/streams/:id/toggle-enabled", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/streams/"+id.String()+"/toggle-enabled", nil)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
