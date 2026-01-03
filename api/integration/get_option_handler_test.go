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

func Test_GetOptionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with option data on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			optionID := "opt-1"
			mockOption := &integration.DriverOption{
				ID:   optionID,
				Name: "Option 1",
			}
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetOption(gomock.Any(), id, optionID).
				Return(mockOption, nil)

			handler := getOptionHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/:id/options/:optionID", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/api/integrations/"+id.String()+"/options/"+optionID,
				nil,
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp integrationOptionResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, optionID, resp.ID)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := getOptionHandler{
				commands: nil,
			}
			r := gin.New()
			r.GET("/api/integrations/:id/options/:optionID", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/invalid/options/opt-1", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("returns 404 Not Found when option does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			optionID := "opt-1"
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetOption(gomock.Any(), id, optionID).
				Return(nil, nil)

			handler := getOptionHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/:id/options/:optionID", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/api/integrations/"+id.String()+"/options/"+optionID,
				nil,
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			optionID := "opt-1"
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				GetOption(gomock.Any(), id, optionID).
				Return(nil, expectedErr)

			handler := getOptionHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/:id/options/:optionID", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/api/integrations/"+id.String()+"/options/"+optionID,
				nil,
			)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
