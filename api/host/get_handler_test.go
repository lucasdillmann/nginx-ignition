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

func Test_GetHandler(t *testing.T) {
	id := uuid.New()
	hostRecord := &host.Host{
		ID:          id,
		DomainNames: []string{"test.com"},
	}

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK when host is found", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(hostRecord, nil)

			settingsCommands := settings.NewMockedCommands(ctrl)
			settingsCommands.EXPECT().
				Get(gomock.Any()).
				Return(&settings.Settings{}, nil)

			router := gin.New()
			handler := getHandler{
				hostCommands:     commands,
				settingsCommands: settingsCommands,
			}
			router.GET("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var response hostResponseDTO
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, hostRecord.DomainNames, response.DomainNames)
		})

		t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
			router := gin.New()
			handler := getHandler{
				hostCommands:     nil,
				settingsCommands: nil,
			}
			router.GET("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/hosts/invalid-uuid", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("returns 404 Not Found when record does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(nil, nil)

			router := gin.New()
			handler := getHandler{
				hostCommands:     commands,
				settingsCommands: nil,
			}
			router.GET("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			commands := host.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			router := gin.New()
			handler := getHandler{
				hostCommands:     commands,
				settingsCommands: nil,
			}
			router.GET("/api/hosts/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/hosts/"+id.String(), nil)

			assert.PanicsWithValue(t, expectedErr, func() {
				router.ServeHTTP(w, req)
			})
		})
	})
}
