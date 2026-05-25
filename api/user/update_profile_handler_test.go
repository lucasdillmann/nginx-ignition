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

func Test_updateProfileHandler(t *testing.T) {
	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			id := uuid.New()
			payload := userProfileUpdateRequestDTO{
				Name:     new("Updated Name"),
				Username: new("updateduser"),
			}

			commands := user.NewMockedCommands(controller)
			commands.EXPECT().
				UpdateProfile(gomock.Any(), id, *payload.Name, *payload.Username).
				Return(nil)

			handler := updateProfileHandler{
				commands: commands,
			}
			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
				ginContext.Next()
			})
			engine.PUT("/current", handler.handle)

			body, _ := json.Marshal(payload)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodPut,
				"/current",
				bytes.NewBuffer(body),
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})
	})
}
