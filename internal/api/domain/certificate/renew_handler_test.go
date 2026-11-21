package certificate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func Test_renewHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("streams a successful result", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Renew(gomock.Any(), id).
				Return(nil)

			handler := renewHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/certificates/:id/renew", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/certificates/"+id.String()+"/renew", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, "text/event-stream", recorder.Header().Get("Content-Type"))

			frame := lastFrame(t, parseSSE(t, recorder.Body.Bytes()))
			assert.Equal(t, "result", frame.Event)
			assert.Equal(t, true, frame.Data["success"])
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := renewHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.POST("/api/certificates/:id/renew", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/certificates/invalid/renew", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("streams a failed result on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Renew(gomock.Any(), id).
				Return(expectedErr)

			handler := renewHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/certificates/:id/renew", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/certificates/"+id.String()+"/renew", nil)
			engine.ServeHTTP(recorder, request)

			frame := lastFrame(t, parseSSE(t, recorder.Body.Bytes()))
			assert.Equal(t, "result", frame.Event)
			assert.Equal(t, false, frame.Data["success"])
			assert.Equal(t, expectedErr.Error(), frame.Data["errorReason"])
		})
	})
}
