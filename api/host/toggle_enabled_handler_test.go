package host

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/host"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_toggleEnabledHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			hostData := newHost()
			hostData.Enabled = true
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), hostData.ID).
				Return(hostData, nil)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/hosts/:id/toggle-enabled", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"POST",
				"/api/hosts/"+hostData.ID.String()+"/toggle-enabled",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
			assert.False(t, hostData.Enabled)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := toggleEnabledHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.POST("/api/hosts/:id/toggle-enabled", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/hosts/invalid/toggle-enabled", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			expectedErr := errors.New("toggle error")
			commands := host.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			handler := toggleEnabledHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/hosts/:id/toggle-enabled", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/hosts/"+id.String()+"/toggle-enabled", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
