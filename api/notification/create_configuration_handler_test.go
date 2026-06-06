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

func Test_createConfigurationHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 201 Created on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
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
				DoAndReturn(func(_ any, _ uuid.UUID, configuration *notification.Configuration) error {
					assert.Equal(t, userID, configuration.UserID)
					assert.Equal(t, payload.Name, configuration.Name)
					assert.Nil(t, configuration.Categories)
					return nil
				})

			handler := createConfigurationHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.POST("", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusCreated, recorder.Code)
		})
	})
}
