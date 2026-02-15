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

func Test_updatePasswordHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := userPasswordUpdateRequestDTO{
				CurrentPassword: new("oldpassword"),
				NewPassword:     new("newpassword"),
			}

			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				UpdatePassword(gomock.Any(), id, *payload.CurrentPassword, *payload.NewPassword).
				Return(nil)

			handler := updatePasswordHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.PATCH("/current/update-password", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PATCH",
				"/current/update-password",
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updatePasswordHandler{
				commands: nil,
			}
			engine := gin.New()
			engine.PATCH("/current/update-password", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"PATCH",
				"/api/users/invalid/password",
				bytes.NewBufferString("{}"),
			)
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set(
					"ABAC:Subject",
					&authorization.Subject{User: &user.User{ID: uuid.New()}},
				)
				ginContext.Next()
			})
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	})
}
