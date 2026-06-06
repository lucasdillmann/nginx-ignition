package notification

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
	"dillmann.com.br/nginx-ignition/core/notification"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_unreadCountHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with unread count on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				UnreadCount(gomock.Any(), userID).
				Return(3, nil)

			handler := unreadCountHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/unread-count", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/unread-count", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response unreadCountResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, 3, response.Count)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := unreadCountHandler{commands: nil}
			engine := gin.New()
			engine.GET("/unread-count", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/unread-count", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
