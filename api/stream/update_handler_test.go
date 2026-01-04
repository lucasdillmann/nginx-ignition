package stream

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

	"dillmann.com.br/nginx-ignition/core/stream"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_updateHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := newStreamRequest()
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(nil)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/streams/:id", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/streams/"+id.String(),
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/streams/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PUT",
				"/api/streams/invalid",
				bytes.NewBufferString("{}"),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := newStreamRequest()
			expectedErr := assert.AnError
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.PUT("/api/streams/:id", func(ginContext *gin.Context) {
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
			request := httptest.NewRequest(
				"PUT",
				"/api/streams/"+id.String(),
				bytes.NewBuffer(body),
			)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
