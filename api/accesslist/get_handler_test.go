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

func Test_getHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK when list is found", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			accessList := newAccessList()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, idToGet uuid.UUID) (*accesslist.AccessList, error) {
					assert.Equal(t, accessList.ID, idToGet)
					return accessList, nil
				})

			engine := gin.New()
			handler := getHandler{
				commands: commands,
			}
			engine.GET("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/access-lists/"+accessList.ID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response accessListResponseDTO
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, accessList.Name, response.Name)
		})

		t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
			engine := gin.New()
			handler := getHandler{
				commands: nil,
			}
			engine.GET("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/access-lists/invalid-uuid", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("returns 404 Not Found when record does not exist", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, nil)

			engine := gin.New()
			handler := getHandler{
				commands: commands,
			}
			engine.GET("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/access-lists/"+id.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Get(gomock.Any(), id).
				Return(nil, expectedErr)

			engine := gin.New()
			handler := getHandler{
				commands: commands,
			}
			engine.GET("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/access-lists/"+id.String(), nil)

			assert.PanicsWithValue(t, expectedErr, func() {
				engine.ServeHTTP(recorder, request)
			})
		})
	})
}
