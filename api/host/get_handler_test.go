package host

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK when host is found", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			hostData := newHost()
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), hostData.ID).
				Return(hostData, nil)

			settingsCommands := settings.NewMockedCommands(controller)
			settingsCommands.EXPECT().
				Get(gomock.Any()).
				Return(&settings.Settings{}, nil)

			engine := gin.New()
			handler := getHandler{
				hostCommands:     commands,
				settingsCommands: settingsCommands,
			}
			engine.GET("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/hosts/"+hostData.ID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response hostResponseDTO
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, hostData.DomainNames, response.DomainNames)
		})

		t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
			engine := gin.New()
			handler := getHandler{
				hostCommands:     nil,
				settingsCommands: nil,
			}
			engine.GET("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/hosts/invalid-uuid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found when record does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, nil)

			engine := gin.New()
			handler := getHandler{
				hostCommands:     commands,
				settingsCommands: nil,
			}
			engine.GET("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			engine := gin.New()
			handler := getHandler{
				hostCommands:     commands,
				settingsCommands: nil,
			}
			engine.GET("/api/hosts/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)

			assert.PanicsWithValue(t, expectedErr, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
