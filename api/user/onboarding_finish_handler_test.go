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
	"dillmann.com.br/nginx-ignition/core/user"
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
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(false, nil)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			commands.EXPECT().
				Authenticate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(user.AuthenticationSuccessful, &user.User{
					ID:       uuid.New(),
					Username: "admin",
				}, nil)

			cfg := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.security.jwt.secret": "1234567890123456789012345678901234567890123456789012345678901234",
			})
			authorizer, _ := authorization.New(cfg, commands)
			handler := onboardingFinishHandler{
				commands:   commands,
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

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newUserRequest()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				OnboardingCompleted(gomock.Any()).
				Return(false, nil)

			expectedErr := assert.AnError
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(expectedErr)

			cfg := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.security.jwt.secret": "1234567890123456789012345678901234567890123456789012345678901234",
			})
			authorizer, _ := authorization.New(cfg, commands)
			handler := onboardingFinishHandler{
				commands:   commands,
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
