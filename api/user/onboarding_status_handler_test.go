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

func Test_OnboardingStatusHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with onboarding status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := user.NewMockedCommands(ctrl)
		commands.EXPECT().
			OnboardingCompleted(gomock.Any()).
			Return(true, nil)

		handler := onboardingStatusHandler{commands}
		r := gin.New()
		r.GET("/api/users/onboarding/status", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/users/onboarding/status", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp userOnboardingStatusResponseDTO
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Finished)
	})
}
