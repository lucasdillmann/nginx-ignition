package certificate

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

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func Test_IssueHandler_Handle(t *testing.T) {
	payload := issueCertificateRequest{
		ProviderID:  "test",
		DomainNames: []string{"example.com"},
	}
	body, _ := json.Marshal(payload)

	t.Run("returns 200 OK with success flag on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		mockCert := &certificate.Certificate{ID: id}

		commands := certificate.NewMockedCommands(ctrl)
		commands.EXPECT().
			Issue(gomock.Any(), gomock.Any()).
			Return(mockCert, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/certificates", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := issueHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp issueCertificateResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
		assert.Equal(t, &id, resp.CertificateID)
	})

	t.Run("returns 200 OK with error reasoning on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := certificate.NewMockedCommands(ctrl)
		commands.EXPECT().
			Issue(gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/certificates", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := issueHandler{commands}
		handler.handle(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp issueCertificateResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
		assert.Equal(t, expectedErr.Error(), *resp.ErrorReason)
	})
}
