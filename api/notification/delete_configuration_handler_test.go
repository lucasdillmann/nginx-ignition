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

func Test_deleteConfigurationHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				DeleteConfiguration(gomock.Any(), userID, configurationID).
				Return(nil)

			handler := deleteConfigurationHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.DELETE("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, "/"+configurationID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := deleteConfigurationHandler{commands: nil}
			engine := gin.New()
			engine.DELETE("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, "/"+uuid.New().String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			userID := uuid.New()
			handler := deleteConfigurationHandler{commands: nil}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.DELETE("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, "/invalid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
