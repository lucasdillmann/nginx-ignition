package accesslist

import (
	"context"
	"encoding/json"
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

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_GetHandler_Handle(t *testing.T) {
	id := uuid.New()
	accessList := &accesslist.AccessList{
		ID:   id,
		Name: "Test List",
	}

	t.Run("returns 200 OK when list is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, getID uuid.UUID) (*accesslist.AccessList, error) {
				assert.Equal(t, id, getID)
				return accessList, nil
			})

		router := gin.New()
		router.GET("/api/access-lists/:id", getHandler{commands}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/access-lists/"+id.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response accessListResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, accessList.Name, response.Name)
	})

	t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
		router := gin.New()
		router.GET("/api/access-lists/:id", getHandler{nil}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/access-lists/invalid-uuid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 Not Found when record does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Return(nil, nil)

		router := gin.New()
		router.GET("/api/access-lists/:id", getHandler{commands}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/access-lists/"+id.String(), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("panics when command returns error", func(t *testing.T) {
		expectedErr := errors.New("command error")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		commands := accesslist.NewMockedCommands(ctrl)
		commands.EXPECT().
			Get(gomock.Any(), gomock.Any()).
			Return(nil, expectedErr)

		router := gin.New()
		router.GET("/api/access-lists/:id", getHandler{commands}.handle)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/access-lists/"+id.String(), nil)

		assert.PanicsWithValue(t, expectedErr, func() {
			router.ServeHTTP(w, req)
		})
	})
}
