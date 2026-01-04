package stream

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_listHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with stream list on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			page := newStreamPage()
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(page, nil)

			handler := listHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/streams", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/streams?pageSize=10&pageNumber=1", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response pagination.PageDTO[streamResponseDTO]
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Len(t, response.Contents, 1)
		})

		t.Run("panics on command error", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			expectedErr := assert.AnError
			commands := stream.NewMockedCommands(controller)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			handler := listHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.GET("/api/streams", func(ginContext *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(ginContext)
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/streams", nil)

			assert.Panics(t, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
