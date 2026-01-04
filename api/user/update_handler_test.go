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

func Test_UpdateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := uuid.New()
			payload := &userRequestDTO{
				Name:     ptr.Of("Updated Name"),
				Username: ptr.Of("updateduser"),
				Enabled:  ptr.Of(true),
			}

			commands := user.NewMockedCommands(ctrl)
			commands.EXPECT().
				Save(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)

			handler := updateHandler{
				commands: commands,
			}
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: uuid.New()}})
				c.Next()
			})
			r.PUT("/api/users/:id", handler.handle)

			jsonPayload, _ := json.Marshal(payload)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT",
				"/api/users/"+id.String(),
				bytes.NewBuffer(jsonPayload),
			)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNoContent, w.Code)
		})

		t.Run("returns 404 Not Found on invalid ID", func(t *testing.T) {
			handler := updateHandler{
				commands: nil,
			}
			r := gin.New()
			r.PUT("/api/users/:id", handler.handle)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/users/invalid", bytes.NewBufferString("{}"))
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	})
}
