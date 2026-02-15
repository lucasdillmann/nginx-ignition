package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_totpDisableHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				DisableTOTP(gomock.Any(), id).
				Return(nil)

			handler := totpDisableHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.POST("/current/totp/disable", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/current/totp/disable", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})
	})
}
