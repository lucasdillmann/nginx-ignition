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
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
	"dillmann.com.br/nginx-ignition/core/stream"
)

func Test_ListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 200 OK with stream list on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStreams := []stream.Stream{
				{
					Name: "stream-1",
				},
				{
					Name: "stream-2",
				},
			}
			page := &corepagination.Page[stream.Stream]{
				Contents: mockStreams,
			}

			commands := stream.NewMockedCommands(ctrl)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(page, nil)

			handler := listHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/streams", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/streams?pageSize=10&pageNumber=1", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var resp pagination.PageDTO[streamResponseDTO]
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Len(t, resp.Contents, 2)
		})

		t.Run("panics on command error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expectedErr := assert.AnError
			commands := stream.NewMockedCommands(ctrl)
			commands.EXPECT().
				List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, expectedErr)

			handler := listHandler{
				commands: commands,
			}
			r := gin.New()
			r.GET("/api/streams", func(c *gin.Context) {
				defer func() {
					if r := recover(); r != nil {
						assert.Equal(t, expectedErr, r)
						panic(r)
					}
				}()
				handler.handle(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/streams", nil)

			assert.Panics(t, func() {
				r.ServeHTTP(w, req)
			})
		})
	})
}
