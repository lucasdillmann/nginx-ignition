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

func Test_getConfigurationHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with configuration data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			configuration := sampleConfiguration(userID)
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				GetConfiguration(gomock.Any(), userID, configuration.ID).
				Return(&configuration, nil)

			handler := getConfigurationHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+configuration.ID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response configurationResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, configuration.ID, response.ID)
			assert.Equal(t, configuration.Name, response.Name)
		})

		t.Run("returns 401 Unauthorized without subject", func(t *testing.T) {
			handler := getConfigurationHandler{commands: nil}
			engine := gin.New()
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+uuid.New().String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			userID := uuid.New()
			handler := getConfigurationHandler{commands: nil}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/invalid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found when configuration does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			userID := uuid.New()
			configurationID := uuid.New()
			commands := notification.NewMockedCommands(controller)
			commands.EXPECT().
				GetConfiguration(gomock.Any(), userID, configurationID).
				Return(nil, nil)

			handler := getConfigurationHandler{commands: commands}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: userID}})
				ginContext.Next()
			})
			engine.GET("/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/"+configurationID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
