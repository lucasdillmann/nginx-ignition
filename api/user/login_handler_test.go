package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_loginHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with token on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := userLoginRequestDTO{
				Username: ptr.Of("johndoe"),
				Password: ptr.Of("password"),
			}

			subject := newUser()
			subject.Username = "johndoe"
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				Authenticate(gomock.Any(), *payload.Username, *payload.Password).
				Return(subject, nil)

			cfg := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.security.jwt.secret": "1234567890123456789012345678901234567890123456789012345678901234",
			})
			authorizer, _ := authorization.New(cfg, commands)
			handler := loginHandler{
				commands:   commands,
				authorizer: authorizer,
			}
			engine := gin.New()
			engine.POST("/api/users/login", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response userLoginResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NotEmpty(t, response.Token)
		})

		t.Run("returns 401 Unauthorized on authentication failure", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := userLoginRequestDTO{
				Username: ptr.Of("johndoe"),
				Password: ptr.Of("password"),
			}

			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				Authenticate(gomock.Any(), *payload.Username, *payload.Password).
				Return(nil, nil)

			cfg := configuration.NewWithOverrides(map[string]string{
				"nginx-ignition.security.jwt.secret": "1234567890123456789012345678901234567890123456789012345678901234",
			})
			authorizer, _ := authorization.New(cfg, commands)
			handler := loginHandler{
				commands:   commands,
				authorizer: authorizer,
			}
			engine := gin.New()
			engine.POST("/api/users/login", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
