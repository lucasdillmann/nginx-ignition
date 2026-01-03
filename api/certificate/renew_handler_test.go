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

func Test_RenewHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with success flag on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		commands := certificate.NewMockedCommands(ctrl)
		commands.EXPECT().
			Renew(gomock.Any(), id).
			Return(nil)

		handler := renewHandler{commands}
		r := gin.New()
		r.POST("/api/certificates/:id/renew", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/certificates/"+id.String()+"/renew", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp renewCertificateResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := renewHandler{nil}
		r := gin.New()
		r.POST("/api/certificates/:id/renew", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/certificates/invalid/renew", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 200 OK with error reasoning on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		expectedErr := assert.AnError
		commands := certificate.NewMockedCommands(ctrl)
		commands.EXPECT().
			Renew(gomock.Any(), id).
			Return(expectedErr)

		handler := renewHandler{commands}
		r := gin.New()
		r.POST("/api/certificates/:id/renew", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/certificates/"+id.String()+"/renew", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp renewCertificateResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, expectedErr.Error(), *resp.ErrorReason)
	})
}
