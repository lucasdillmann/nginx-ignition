package vpn

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_DeleteHandler_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Delete(gomock.Any(), id).
			Return(nil)

		handler := deleteHandler{commands}
		r := gin.New()
		r.DELETE("/api/vpns/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/vpns/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
		handler := deleteHandler{nil}
		r := gin.New()
		r.DELETE("/api/vpns/:id", handler.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/vpns/invalid", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics on command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uuid.New()
		expectedErr := assert.AnError
		commands := vpn.NewMockedCommands(ctrl)
		commands.EXPECT().
			Delete(gomock.Any(), id).
			Return(expectedErr)

		handler := deleteHandler{commands}
		r := gin.New()
		r.DELETE("/api/vpns/:id", func(c *gin.Context) {
			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, expectedErr, r)
					panic(r)
				}
			}()
			handler.handle(c)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/vpns/"+id.String(), nil)

		assert.Panics(t, func() {
			r.ServeHTTP(w, req)
		})
	})
}
