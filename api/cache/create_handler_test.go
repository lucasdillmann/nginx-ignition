package cache

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
)

func Test_CreateHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := cacheRequestDTO{
		Name:            "New Cache",
		StoragePath:     ptr.Of("/var/lib/nginx/cache"),
		InactiveSeconds: ptr.Of(3600),
	}
	body, _ := json.Marshal(payload)

	t.Run("returns 201 Created on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := cache.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := createHandler{commands}
		r := gin.New()
		r.POST("/api/caches", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/caches", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp cacheResponseDTO
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, payload.Name, resp.Name)
		assert.NotEqual(t, "", resp.ID.String())
	})

	t.Run("panics on invalid JSON", func(t *testing.T) {
		handler := createHandler{nil}
		r := gin.New()
		r.POST("/api/caches", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/caches", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := cache.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		handler := createHandler{commands}
		r := gin.New()
		r.POST("/api/caches", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/caches", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
