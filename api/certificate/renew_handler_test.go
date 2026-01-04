package certificate

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_renewHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with success flag on success", func(t *testing.T) {
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
			var response renewCertificateResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.True(t, response.Success)
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

		t.Run("returns 200 OK with error reasoning on command error", func(t *testing.T) {
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

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response renewCertificateResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.False(t, response.Success)
			assert.Equal(t, expectedErr.Error(), *response.ErrorReason)
		})
	})
}
