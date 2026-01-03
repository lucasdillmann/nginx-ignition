package host

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/host"
)

func Test_CreateHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	payload := &hostRequestDTO{
		Enabled:     ptr.Of(true),
		DomainNames: []string{"example.com"},
		FeatureSet:  &featureSetDTO{},
	}
	body, _ := json.Marshal(payload)

	t.Run("returns 201 Created on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := host.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := createHandler{commands}
		r := gin.New()
		r.POST("/api/hosts", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/hosts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotNil(t, resp["id"])
	})

	t.Run("panics on invalid JSON", func(t *testing.T) {
		handler := createHandler{nil}
		r := gin.New()
		r.POST("/api/hosts", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/hosts", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := host.NewMockedCommands(ctrl)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		handler := createHandler{commands}
		r := gin.New()
		r.POST("/api/hosts", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/hosts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
