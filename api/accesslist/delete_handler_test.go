package accesslist

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func Test_DeleteHandler_Handle(t *testing.T) {
	id := uuid.New()

	t.Run("returns 204 No Content on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, delID uuid.UUID) error {
				assert.Equal(t, id, delID)
				return nil
			})

		router := gin.New()
		router.DELETE("/api/access-lists/:id", deleteHandler{commands}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/access-lists/"+id.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/api/access-lists/:id", deleteHandler{nil}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/access-lists/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		expectedErr := errors.New("command error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Delete(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		router := gin.New()
		router.DELETE("/api/access-lists/:id", deleteHandler{commands}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/access-lists/"+id.String(), nil)

		assert.PanicsWithValue(t, expectedErr, func() {
			router.ServeHTTP(w, req)
		})
	})
}
