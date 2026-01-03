package host

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/host"
)

func Test_ToggleEnabledHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		mockHost := &host.Host{ID: id, Enabled: true}
		commands := host.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), id).
			Return(mockHost, nil)
		commands.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		handler := toggleEnabledHandler{commands}
		r := gin.New()
		r.POST("/api/hosts/:id/toggle-enabled", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/hosts/"+id.String()+"/toggle-enabled", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.False(t, mockHost.Enabled)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := toggleEnabledHandler{nil}
		r := gin.New()
		r.POST("/api/hosts/:id/toggle-enabled", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/hosts/invalid/toggle-enabled", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		expectedErr := errors.New("toggle error")
		commands := host.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), id).
			Return(nil, expectedErr)

		handler := toggleEnabledHandler{commands}
		r := gin.New()
		r.POST("/api/hosts/:id/toggle-enabled", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/hosts/"+id.String()+"/toggle-enabled", nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
