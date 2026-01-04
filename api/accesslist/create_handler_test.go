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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_createHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 201 Created on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			payload := newAccessListRequestDTO()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, accessList *accesslist.AccessList) error {
					assert.Equal(t, *payload.Name, accessList.Name)
					return nil
				})

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/access-lists",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := createHandler{
				commands: commands,
			}
			handler.handle(ginContext)

			assert.Equal(t, http.StatusCreated, recorder.Code)
			var response map[string]string
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NotEmpty(t, response["id"])
		})

		t.Run("panics on invalid JSON", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/access-lists",
				bytes.NewBufferString("invalid json"),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := createHandler{
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

			payload := newAccessListRequestDTO()
			commands := accesslist.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any()).
				Return(expectedErr)

			recorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(recorder)
			body, _ := json.Marshal(payload)
			ginContext.Request = httptest.NewRequest(
				"POST",
				"/api/access-lists",
				bytes.NewBuffer(body),
			)
			ginContext.Request.Header.Set("Content-Type", "application/json")

			handler := createHandler{
				commands: commands,
			}
			assert.PanicsWithValue(t, expectedErr, func() {
				handler.handle(ginContext)
			})
		})
	})
}
