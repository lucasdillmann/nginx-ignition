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

func Test_getOptionHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with option data on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			optionID := "opt-1"
			option := &integration.DriverOption{
				ID:   optionID,
				Name: "Option 1",
			}
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				GetOption(gomock.Any(), id, optionID).
				Return(option, nil)

			handler := getOptionHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id/options/:optionID", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/integrations/"+id.String()+"/options/"+optionID,
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response integrationOptionResponse
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Equal(t, optionID, response.ID)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := getOptionHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id/options/:optionID", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/integrations/invalid/options/opt-1", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found when option does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			optionID := "opt-1"
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				GetOption(gomock.Any(), id, optionID).
				Return(nil, nil)

			handler := getOptionHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id/options/:optionID", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/integrations/"+id.String()+"/options/"+optionID,
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			optionID := "opt-1"
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				GetOption(gomock.Any(), id, optionID).
				Return(nil, expectedErr)

			handler := getOptionHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/integrations/:id/options/:optionID", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/api/integrations/"+id.String()+"/options/"+optionID,
				nil,
			)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
