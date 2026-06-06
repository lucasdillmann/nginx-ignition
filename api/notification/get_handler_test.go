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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with notification data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			subject := sampleNotification(userID)
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				GetNotification(gomock.Any(), userID, subject.ID).
				Return(subject, nil)

			handler := getHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+subject.ID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response notificationResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, subject.ID, response.ID)
			assert.Equal(t, subject.Title, response.Title)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := getHandler{commands: nil}
			engine := gin.New()
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+uuid.New().String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("returns 404 Not Found when notification does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			notificationID := uuid.New()
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				GetNotification(gomock.Any(), userID, notificationID).
				Return(nil, nil)

			handler := getHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+notificationID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
