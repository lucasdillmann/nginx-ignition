package accesslist

import (
	"bytes"
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

func Test_updateHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := newAccessListRequestDTO()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, accessList *accesslist.AccessList) error {
					assert.Equal(t, id, accessList.ID)
					assert.Equal(t, *payload.Name, accessList.Name)
					return nil
				})

			engine := gin.New()
			handler := updateHandler{
				commands: commands,
			}
			engine.PUT("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			body, _ := json.Marshal(payload)
			request := httptest.NewRequest(
				"PUT",
				"/api/access-lists/"+id.String(),
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found when ID is invalid", func(t *testing.T) {
			payload := newAccessListRequestDTO()
			engine := gin.New()
			handler := updateHandler{
				commands: nil,
			}
			engine.PUT("/api/access-lists/:id", handler.handle)

			recorder := httptest.NewRecorder()
			body, _ := json.Marshal(payload)
			request := httptest.NewRequest(
				"PUT",
				"/api/access-lists/invalid-uuid",
				bytes.NewBuffer(body),
			)
			request.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			id := uuid.New()
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Params = gin.Params{
				{
					Key:   "id",
					Value: id.String(),
				},
			}
			ginContext.Request = httptest.NewRequest(
				"PUT",
				"/api/access-lists/"+id.String(),
				bytes.NewBufferString("invalid json"),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := updateHandler{
				commands: nil,
			}
			assert.Panics(t, func() {
				handler.handle(ginContext)
			})
		})

		t.Run("panics when command returns error", func(t *testing.T) {
			expectedErr := errors.New("command error")
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := newAccessListRequestDTO()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Params = gin.Params{
				{
					Key:   "id",
					Value: id.String(),
				},
			}
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"PUT",
				"/api/access-lists/"+id.String(),
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := updateHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
