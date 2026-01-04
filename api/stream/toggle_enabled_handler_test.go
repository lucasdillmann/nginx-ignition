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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_toggleEnabledHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			subject := newStream()
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), subject.ID).
				Return(subject, nil)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PATCH("/api/streams/:id/toggle-enabled", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PATCH",
				"/api/streams/"+subject.ID.String()+"/toggle-enabled",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
			assert.False(t, subject.Enabled)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := toggleEnabledHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PATCH("/api/streams/:id/toggle-enabled", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PATCH", "/api/streams/invalid/toggle-enabled", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found when stream does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, nil)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PATCH("/api/streams/:id/toggle-enabled", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PATCH",
				"/api/streams/"+id.String()+"/toggle-enabled",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PATCH("/api/streams/:id/toggle-enabled", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PATCH",
				"/api/streams/"+id.String()+"/toggle-enabled",
				nil,
			)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
