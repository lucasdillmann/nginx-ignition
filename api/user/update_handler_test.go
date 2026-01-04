package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
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
			payload := newUserRequest()
			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)

			handler := updateHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set(
					"ABAC:Subject",
					&authorization.Subject{User: &user.User{ID: uuid.New()}},
				)
				ginContext.Next()
			})
			engine.PUT("/api/users/:id", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PUT", "/api/users/"+id.String(), bytes.NewBuffer(body))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PUT("/api/users/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PUT", "/api/users/invalid", bytes.NewBufferString("{}"))
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
