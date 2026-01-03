package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_LogoutHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Setenv(
			"NGINX_IGNITION_SECURITY_JWT_SECRET",
			"1234567890123456789012345678901234567890123456789012345678901234",
		)
		commands := user.NewMockedCommands(ctrl)
		authorizer, _ := authorization.New(configuration.New(), commands)

		handler := logoutHandler{authorizer}
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("ABAC:Subject", &authorization.Subject{TokenID: "token-id"})
			c.Next()
		})
		r.POST("/api/users/logout", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/users/logout", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
