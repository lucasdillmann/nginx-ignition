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

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_onboardingFinishHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with token on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newUserRequest()
			userCommands := user.NewMockedCommands(controller)
			userCommands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(false, nil)
			userCommands.EXPECT().
				FinishOnboarding(gomock.Any(), gomock.Any()).
				Return(nil)
			userCommands.EXPECT().
				Authenticate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(user.AuthenticationSuccessful, &user.User{
					ID:       uuid.New(),
					Username: "admin",
				}, nil)

			authorizer, _ := authorization.New(
				configuration.New(),
				userCommands,
				newAuthorizationCommands(controller),
			)
			handler := onboardingFinishHandler{
				commands:   userCommands,
				authorizer: authorizer,
			}
			engine := gin.New()
			engine.POST("/api/users/onboarding/finish", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/users/onboarding/finish",
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("returns 403 when onboarding already completed", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newUserRequest()
			userCommands := user.NewMockedCommands(controller)
			userCommands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(true, nil)

			authorizer, _ := authorization.New(
				configuration.New(),
				userCommands,
				newAuthorizationCommands(controller),
			)
			handler := onboardingFinishHandler{
				commands:   userCommands,
				authorizer: authorizer,
			}
			engine := gin.New()
			engine.POST("/api/users/onboarding/finish", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/users/onboarding/finish",
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("returns 403 when FinishOnboarding reports already completed", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newUserRequest()
			userCommands := user.NewMockedCommands(controller)
			userCommands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(false, nil)
			userCommands.EXPECT().
				FinishOnboarding(gomock.Any(), gomock.Any()).
				Return(user.ErrOnboardingAlreadyCompleted)

			authorizer, _ := authorization.New(
				configuration.New(),
				userCommands,
				newAuthorizationCommands(controller),
			)
			handler := onboardingFinishHandler{
				commands:   userCommands,
				authorizer: authorizer,
			}
			engine := gin.New()
			engine.POST("/api/users/onboarding/finish", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/users/onboarding/finish",
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newUserRequest()
			userCommands := user.NewMockedCommands(controller)

			expectedErr := assert.AnError
			userCommands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(false, nil)
			userCommands.EXPECT().
				FinishOnboarding(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			authorizer, _ := authorization.New(
				configuration.New(),
				userCommands,
				newAuthorizationCommands(controller),
			)
			handler := onboardingFinishHandler{
				commands:   userCommands,
				authorizer: authorizer,
			}
			engine := gin.New()
			engine.POST("/api/users/onboarding/finish", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/users/onboarding/finish",
				bytes.NewBuffer(body),
			)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
