package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/integration"
)

func Test_AvailableDriversHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with available drivers", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDrivers := []integration.AvailableDriver{
				{
					ID:   "docker",
					Name: "Docker",
				},
				{
					ID:   "swarm",
					Name: "Swarm",
				},
			}

			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetAvailableDrivers(gomock.Any()).
				Return(mockDrivers, nil)

			handler := availableDriversHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/available-drivers", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/available-drivers", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp []integrationDriverResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Len(t, resp, 2)
			assert.Equal(t, "docker", resp[0].ID)
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetAvailableDrivers(gomock.Any()).
				Return(nil, expectedErr)

			handler := availableDriversHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/available-drivers", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/available-drivers", nil)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
