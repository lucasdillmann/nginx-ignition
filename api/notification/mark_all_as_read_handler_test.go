package notification

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/notification"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_markAllAsReadHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				MarkAllAsRead(gomock.Any(), userID).
				Return(nil)

			handler := markAllAsReadHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.POST("/mark-all-as-read", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/mark-all-as-read", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := markAllAsReadHandler{commands: nil}
			engine := gin.New()
			engine.POST("/mark-all-as-read", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/mark-all-as-read", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
