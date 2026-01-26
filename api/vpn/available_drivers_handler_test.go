package vpn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_availableDriversHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with available drivers", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			drivers := []vpn.AvailableDriver{
				{
					ID:   "openvpn",
					Name: i18n.Static("OpenVPN"),
				},
				{
					ID:   "wireguard",
					Name: i18n.Static("WireGuard"),
				},
			}

			commands := vpn.NewMockedCommands(controller)
			commands.EXPECT().
				GetAvailableDrivers(gomock.Any()).
				Return(drivers, nil)

			handler := availableDriversHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/vpns/available-drivers", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/vpns/available-drivers", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []vpnDriverResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response, 2)
			assert.Equal(t, "openvpn", response[0].ID)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := vpn.NewMockedCommands(controller)
			commands.EXPECT().
				GetAvailableDrivers(gomock.Any()).
				Return(nil, expectedErr)

			handler := availableDriversHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/vpns/available-drivers", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/vpns/available-drivers", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
