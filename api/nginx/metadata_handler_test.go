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

func Test_MetadataHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with metadata", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockMetadata := &nginx.Metadata{
			Version: "1.21.0",
		}
		commands := nginx.NewMockedCommands(ctrl)
		commands.EXPECT().
			GetMetadata(gomock.Any()).
			Return(mockMetadata, nil)

		handler := metadataHandler{commands}
		r := gin.New()
		r.GET("/api/nginx/metadata", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/nginx/metadata", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "1.21.0", resp["version"])
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := nginx.NewMockedCommands(ctrl)
		commands.EXPECT().
			GetMetadata(gomock.Any()).
			Return(nil, expectedErr)

		handler := metadataHandler{commands}
		r := gin.New()
		r.GET("/api/nginx/metadata", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/nginx/metadata", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
