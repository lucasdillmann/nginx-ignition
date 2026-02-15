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

func Test_totpStatusHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with enabled true when TOTP is activated", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				GetTOTPStatus(gomock.Any(), id).
				Return(true, nil)

			handler := totpStatusHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.GET("/current/totp/status", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/totp/status", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response totpStatusResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.True(t, response.Enabled)
		})

		t.Run("returns 200 OK with enabled false when TOTP is not activated", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				GetTOTPStatus(gomock.Any(), id).
				Return(false, nil)

			handler := totpStatusHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.GET("/current/totp/status", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/totp/status", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response totpStatusResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.False(t, response.Enabled)
		})
	})
}
