package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func Test_apiTokenDeleteHandler(t *testing.T) {
	setup := func(t *testing.T, id uuid.UUID) (*user.MockedCommands, *gin.Engine) {
		controller := gomock.NewController(t)
		commands := user.NewMockedCommands(controller)
		handler := apiTokenDeleteHandler{commands: commands}

		engine := gin.New()
		engine.Use(func(ginContext *gin.Context) {
			ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
			ginContext.Next()
		})
		engine.DELETE("/current/tokens/:id", handler.handle)

		return commands, engine
	}

	t.Run("handle", func(t *testing.T) {
		t.Run("returns 204 No Content on success", func(t *testing.T) {
			usr := newUser()
			tokenID := uuid.New()

			commands, engine := setup(t, usr.ID)
			commands.EXPECT().DeleteAPIToken(gomock.Any(), usr.ID, tokenID).Return(nil)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/current/tokens/"+tokenID.String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusNoContent, recorder.Code)
		})

		t.Run("panics for an invalid token ID", func(t *testing.T) {
			usr := newUser()

			_, engine := setup(t, usr.ID)

			assert.Panics(t, func() {
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest("DELETE", "/current/tokens/not-a-uuid", nil)
				engine.ServeHTTP(recorder, request)
			})
		})

		t.Run("returns 401 Unauthorized without a subject", func(t *testing.T) {
			controller := gomock.NewController(t)
			commands := user.NewMockedCommands(controller)
			handler := apiTokenDeleteHandler{commands: commands}
			engine := gin.New()
			engine.DELETE("/current/tokens/:id", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/current/tokens/"+uuid.New().String(), nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})
	})
}
