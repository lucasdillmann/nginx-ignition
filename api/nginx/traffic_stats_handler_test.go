package nginx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

func Test_trafficStatsHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with traffic stats", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			statsData := newTrafficStats()
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetTrafficStats(gomock.Any()).
				Return(statsData, nil)

			handler := trafficStatsHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/traffic-stats", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/traffic-stats", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response trafficStatsResponseDTO
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, statsData.HostName, response.HostName)
			assert.Equal(t, statsData.Connections.Active, response.Connections.Active)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := nginx.NewMockedCommands(controller)
			commands.EXPECT().
				GetTrafficStats(gomock.Any()).
				Return(nil, expectedErr)

			handler := trafficStatsHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/nginx/traffic-stats", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/nginx/traffic-stats", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
