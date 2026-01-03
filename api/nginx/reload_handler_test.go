package nginx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/nginx"
)

func init() {
	_ = log.Init()
}

func Test_ReloadHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := nginx.NewMockedCommands(ctrl)
		commands.EXPECT().
			Reload(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := reloadHandler{commands}
		r := gin.New()
		r.POST("/api/nginx/reload", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/nginx/reload", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns 424 Failed Dependency on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := assert.AnError
		commands := nginx.NewMockedCommands(ctrl)
		commands.EXPECT().
			Reload(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		handler := reloadHandler{commands}
		r := gin.New()
		r.POST("/api/nginx/reload", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/nginx/reload", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFailedDependency, w.Code)
		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, expectedErr.Error(), resp["message"])
	})
}
