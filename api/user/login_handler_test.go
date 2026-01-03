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
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func init() {
	_ = log.Init()
}

func Test_LoginHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with token on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		payload := &userLoginRequestDTO{
			Username: ptr.Of("johndoe"),
			Password: ptr.Of("password"),
		}

		mockUser := &user.User{Username: "johndoe"}
		commands := user.NewMockedCommands(ctrl)
		commands.EXPECT().
			Authenticate(gomock.Any(), *payload.Username, *payload.Password).
			Return(mockUser, nil)

		t.Setenv(
			"NGINX_IGNITION_SECURITY_JWT_SECRET",
			"1234567890123456789012345678901234567890123456789012345678901234",
		)
		authorizer, _ := authorization.New(configuration.New(), commands)
		handler := loginHandler{commands, authorizer}
		r := gin.New()
		r.POST("/api/users/login", handler.handle)

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonPayload))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp userLoginResponseDTO
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("returns 401 Unauthorized on authentication failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		payload := &userLoginRequestDTO{
			Username: ptr.Of("johndoe"),
			Password: ptr.Of("password"),
		}

		commands := user.NewMockedCommands(ctrl)
		commands.EXPECT().
			Authenticate(gomock.Any(), *payload.Username, *payload.Password).
			Return(nil, nil)

		t.Setenv(
			"NGINX_IGNITION_SECURITY_JWT_SECRET",
			"1234567890123456789012345678901234567890123456789012345678901234",
		)
		authorizer, _ := authorization.New(configuration.New(), commands)
		handler := loginHandler{commands, authorizer}
		r := gin.New()
		r.POST("/api/users/login", handler.handle)

		jsonPayload, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonPayload))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
