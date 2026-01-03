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

func Test_GetHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with integration data on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			mockIntegration := &integration.Integration{
				ID:   id,
				Name: "integration-1",
			}
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(mockIntegration, nil)

			handler := getHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/"+id.String(), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp integrationResponse
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, id, resp.ID)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := getHandler{
				commands: nil,
			}
			r := gin.New()
			r.GET("/api/integrations/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/invalid", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("returns 404 Not Found when integration does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, nil)

			handler := getHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/"+id.String(), nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(ctrl)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			handler := getHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/integrations/:id", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/integrations/"+id.String(), nil)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
