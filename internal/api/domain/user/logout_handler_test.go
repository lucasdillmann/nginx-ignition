package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_logoutHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			sessionID := uuid.New()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				StopSession(gomock.Any(), sessionID).
				Return(nil)

			handler := logoutHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set(
					"ABAC:Subject",
					&authorization.Subject{SessionID: sessionID},
				)
				ginContext.Next()
			})
			engine.POST("/api/users/logout", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/users/logout", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})
	})
}
