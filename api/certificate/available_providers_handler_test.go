package certificate

import (
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

func Test_availableProvidersHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with providers list on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			providers := []certificate.AvailableProvider{}
			commands := certificate.NewMockedCommands(controller)
			commands.EXPECT().
				AvailableProviders(gomock.Any()).
				Return(providers, nil)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest(
				"GET",
				"/api/certificates/available-providers",
				nil,
			)

			handler := availableProvidersHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []availableProviderResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response, 0)
		})
	})
}
