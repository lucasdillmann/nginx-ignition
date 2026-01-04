package certificate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_issueHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with success flag on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			certificateData := newCertificate()
			payload := newIssueCertificateRequest()
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				Return(certificateData, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response issueCertificateResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.True(t, response.Success)
			assert.Equal(t, &certificateData.ID, response.CertificateID)
		})

		t.Run("returns 200 OK with error reasoning on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newIssueCertificateRequest()
			expectedErr := assert.AnError
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				Issue(gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/certificates",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := issueHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response issueCertificateResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.False(t, response.Success)
			assert.Equal(t, expectedErr.Error(), *response.ErrorReason)
		})
	})
}
