package cache

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/pagination"
	"dillmann.com.br/nginx-ignition/core/cache"
	corepagination "dillmann.com.br/nginx-ignition/core/common/pagination"
)

func Test_ListHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 200 OK with cache list on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCaches := []cache.Cache{
			{Name: "Cache 1"},
			{Name: "Cache 2"},
		}

		page := &corepagination.Page[cache.Cache]{
			Contents: mockCaches,
		}

		commands := cache.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(page, nil)

		handler := listHandler{commands}
		r := gin.New()
		r.GET("/api/caches", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/caches?pageSize=10&pageNumber=1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp pagination.PageDTO[cacheResponseDTO]
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp.Contents, 2)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("list error")
		commands := cache.NewMockedCommands(ctrl)
		commands.EXPECT().
			List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		handler := listHandler{commands}
		r := gin.New()
		r.GET("/api/caches", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/caches", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
