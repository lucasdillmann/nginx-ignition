package integration

import (
	"bytes"
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

func Test_createHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 201 Created on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newIntegrationRequest()
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := createHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/integrations", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/integrations", bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusCreated, recorder.Code)
			var response map[string]uuid.UUID
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NotEqual(t, uuid.Nil, response["id"])
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newIntegrationRequest()
			expectedErr := assert.AnError
			commands := integration.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := createHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.POST("/api/integrations", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/integrations", bytes.NewBuffer(body))

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
