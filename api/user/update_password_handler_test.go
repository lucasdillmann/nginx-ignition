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
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Test_UpdatePasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			payload := &userPasswordUpdateRequestDTO{
				CurrentPassword: ptr.Of("oldpassword"),
				NewPassword:     ptr.Of("newpassword"),
			}

			commands := user.NewMockedCommands(ctrl)
			commands.EXPECT().
				UpdatePassword(gomock.Any(), id, *payload.CurrentPassword, *payload.NewPassword).
				Return(nil)

			handler := updatePasswordHandler{
				commands: commands,
			}
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				c.Next()
			})
			r.PATCH("/current/update-password", handler.handle)

			jsonPayload, _ := json.Marshal(payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PATCH",
				"/current/update-password",
				bytes.NewBuffer(jsonPayload),
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updatePasswordHandler{
				commands: nil,
			}
			r := gin.New()
			r.PATCH("/current/update-password", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PATCH",
				"/api/users/invalid/password",
				bytes.NewBufferString("{}"),
			)
			r.Use(func(c *gin.Context) {
				c.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: uuid.New()}})
				c.Next()
			})
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	})
}
