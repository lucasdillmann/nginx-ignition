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

func Test_listConfigurationsHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with configuration list on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			configuration := sampleConfiguration(userID)
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				ListConfigurations(gomock.Any(), userID).
				Return([]notification.Configuration{configuration}, nil)

			handler := listConfigurationsHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []configurationResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response, 1)
			assert.Equal(t, configuration.ID, response[0].ID)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := listConfigurationsHandler{commands: nil}
			engine := gin.New()
			engine.GET("", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}

func sampleConfiguration(userID uuid.UUID) notification.Configuration {
	return notification.Configuration{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       "Primary SMTP",
		Provider:   "SMTP",
		Enabled:    true,
		Parameters: map[string]any{"host": "smtp.example.com"},
	}
}
