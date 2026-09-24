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

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_loginHandler(t *testing.T) {
	setup := func(t *testing.T) (*user2.MockedCommands, *gin.Engine) {
		controller := gomock.NewController(t)
		commands := user2.NewMockedCommands(controller)
		authorizer, _ := authorization.New(configuration.New(), commands)
		handler := loginHandler{
			commands:   commands,
			authorizer: authorizer,
		}

		engine := gin.New()
		engine.POST("/api/users/login", handler.handle)
		return commands, engine
	}

	performRequest := func(engine *gin.Engine, payload any) *httptest.ResponseRecorder {
		body, _ := json.Marshal(payload)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
		engine.ServeHTTP(recorder, request)
		return recorder
	}

	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with token on success", func(t *testing.T) {
			payload := userLoginRequestDTO{
				Username: new("johndoe"),
				Password: new("password"),
			}

			subject := newUser()
			subject.Username = "johndoe"
			commands, engine := setup(t)
			commands.EXPECT().
				Authenticate(gomock.Any(), *payload.Username, *payload.Password, gomock.Any()).
				Return(user2.AuthenticationSuccessful, subject, nil)

			recorder := performRequest(engine, payload)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response userLoginResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NotEmpty(t, response.Token)
		})

		t.Run("returns 401 Unauthorized on authentication failure", func(t *testing.T) {
			payload := userLoginRequestDTO{
				Username: new("johndoe"),
				Password: new("password"),
			}

			commands, engine := setup(t)
			commands.EXPECT().
				Authenticate(gomock.Any(), *payload.Username, *payload.Password, gomock.Any()).
				Return(user2.AuthenticationFailed, nil, nil)

			recorder := performRequest(engine, payload)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("returns 200 OK when TOTP code is provided", func(t *testing.T) {
			payload := userLoginRequestDTO{
				Username: new("johndoe"),
				Password: new("password"),
				TOTP:     new("123456"),
			}

			subject := newUser()
			subject.Username = "johndoe"
			commands, engine := setup(t)
			commands.EXPECT().
				Authenticate(gomock.Any(), *payload.Username, *payload.Password, "123456").
				Return(user2.AuthenticationSuccessful, subject, nil)

			recorder := performRequest(engine, payload)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response userLoginResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NotEmpty(t, response.Token)
		})

		t.Run("returns 401 Unauthorized with MISSING_TOTP reason", func(t *testing.T) {
			payload := userLoginRequestDTO{
				Username: new("johndoe"),
				Password: new("password"),
			}

			commands, engine := setup(t)
			commands.EXPECT().
				Authenticate(gomock.Any(), *payload.Username, *payload.Password, gomock.Any()).
				Return(user2.AuthenticationMissingTOTP, nil, nil)

			recorder := performRequest(engine, payload)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			var response map[string]any
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, string(user2.AuthenticationMissingTOTP), response["reason"])
		})
	})
}
