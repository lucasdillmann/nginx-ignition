package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/integration"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_getHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with integration data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			integrationData := newIntegration()
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), integrationData.ID).
				Return(integrationData, nil)

			handler := getHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/integrations/"+integrationData.ID.String(),
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response integrationResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, integrationData.ID, response.ID)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := getHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/integrations/invalid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found when integration does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, nil)

			handler := getHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/integrations/"+id.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			handler := getHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/integrations/"+id.String(), nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
