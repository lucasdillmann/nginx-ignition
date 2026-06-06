package notification

import (
	"bytes"
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

func Test_putConfigurationHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 404 when configuration is owned by another user", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			payload := configurationRequest{
				Name:       "Primary SMTP",
				Provider:   "SMTP",
				Enabled:    true,
				Parameters: map[string]any{"host": "smtp.example.com"},
				Categories: json.RawMessage("null"),
			}

			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				SaveConfiguration(gomock.Any(), userID, gomock.Any()).
				Return(notification.ErrConfigurationNotFound)

			handler := putConfigurationHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.PUT("/:id", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodPut,
				"/"+configurationID.String(),
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
