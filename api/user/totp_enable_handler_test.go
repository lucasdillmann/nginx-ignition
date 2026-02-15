package user

import (
	"encoding/json"
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

func Test_totpEnableHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with secret when TOTP is not yet activated", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				GetTOTPStatus(gomock.Any(), id).
				Return(false, nil)
			commands.EXPECT().
				EnableTOTP(gomock.Any(), id).
				Return("otpauth://totp/test?secret=ABC123", nil)

			handler := totpEnableHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.POST("/current/totp/enable", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/current/totp/enable", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response totpEnableResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, "otpauth://totp/test?secret=ABC123", response.Secret)
		})

		t.Run("returns 400 Bad Request when TOTP is already activated", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				GetTOTPStatus(gomock.Any(), id).
				Return(true, nil)

			handler := totpEnableHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.POST("/current/totp/enable", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/current/totp/enable", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	})
}
