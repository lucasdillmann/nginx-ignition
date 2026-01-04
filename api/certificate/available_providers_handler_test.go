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

func Test_AvailableProvidersHandler(t *testing.T) {
	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with providers list on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProviders := []certificate.AvailableProvider{}

			commands := certificate.NewMockedCommands(ctrl)
			commands.EXPECT().
				AvailableProviders(gomock.Any()).
				Return(mockProviders, nil)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest("GET", "/api/certificates/available-providers", nil)

			handler := availableProvidersHandler{
				commands: commands,
			}
			handler.handle(ctx)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp []availableProviderResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Len(t, resp, 0)
		})
	})
}
