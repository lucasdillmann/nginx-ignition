package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/user"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_onboardingStatusHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with onboarding status", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(true, nil)

			handler := onboardingStatusHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/users/onboarding/status", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/users/onboarding/status", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response userOnboardingStatusResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.True(t, response.Finished)
		})
	})
}
