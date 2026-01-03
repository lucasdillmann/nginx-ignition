package user

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

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_OnboardingFinishHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		payload := &userRequestDTO{
			Name:     ptr.Of("Admin"),
			Username: ptr.Of("admin"),
			Password: ptr.Of("password"),
		}

		t.Setenv(
			"NGINX_IGNITION_SECURITY_JWT_SECRET",
			"1234567890123456789012345678901234567890123456789012345678901234",
		)
		commands := user.NewMockedCommands(ctrl)
		commands.EXPECT().
			OnboardingCompleted(gomock.Any()).
			Return(false, nil)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		commands.EXPECT().
			Authenticate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&user.User{ID: uuid.New(), Username: "admin"}, nil)

		authorizer, _ := authorization.New(configuration.New(), commands)
		handler := onboardingFinishHandler{commands, authorizer}
		r := gin.New()
		r.POST("/api/users/onboarding/finish", handler.handle)

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			"POST",
			"/api/users/onboarding/finish",
			bytes.NewBuffer(jsonPayload),
		)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		payload := &userRequestDTO{
			Name:     ptr.Of("Admin"),
			Username: ptr.Of("admin"),
		}

		t.Setenv(
			"NGINX_IGNITION_SECURITY_JWT_SECRET",
			"1234567890123456789012345678901234567890123456789012345678901234",
		)
		commands := user.NewMockedCommands(ctrl)
		commands.EXPECT().
			OnboardingCompleted(gomock.Any()).
			Return(false, nil)

		expectedErr := assert.AnError
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedErr)

		authorizer, _ := authorization.New(configuration.New(), commands)
		handler := onboardingFinishHandler{commands, authorizer}
		r := gin.New()
		r.POST("/api/users/onboarding/finish", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			"POST",
			"/api/users/onboarding/finish",
			bytes.NewBuffer(jsonPayload),
		)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
