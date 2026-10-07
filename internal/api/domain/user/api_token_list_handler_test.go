package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func Test_apiTokenListHandler(t *testing.T) {
	setup := func(t *testing.T, id uuid.UUID) (*user.MockedCommands, *gin.Engine) {
		controller := gomock.NewController(t)
		commands := user.NewMockedCommands(controller)
		handler := apiTokenListHandler{commands: commands}

		engine := gin.New()
		engine.Use(func(ginContext *gin.Context) {
			ginContext.Set("ABAC:Subject", &authorization.Subject{User: &user.User{ID: id}})
			ginContext.Next()
		})
		engine.GET("/current/tokens", handler.handle)

		return commands, engine
	}

	t.Run("handle", func(t *testing.T) {
		t.Run("returns 200 OK with the user tokens", func(t *testing.T) {
			usr := newUser()
			token := newAPIToken(usr)

			commands, engine := setup(t, usr.ID)
			commands.EXPECT().
				ListAPITokens(gomock.Any(), usr.ID).
				Return([]user.APIToken{*token}, nil)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/tokens", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []apiTokenResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			require.Len(t, response, 1)
			assert.Equal(t, token.ID, response[0].ID)
			assert.Equal(t, token.Name, response[0].Name)
			assert.Equal(t, token.Expiration, response[0].Expiration)
		})

		t.Run("returns 200 OK with an empty list", func(t *testing.T) {
			usr := newUser()

			commands, engine := setup(t, usr.ID)
			commands.EXPECT().ListAPITokens(gomock.Any(), usr.ID).Return(nil, nil)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/tokens", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response []apiTokenResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Empty(t, response)
		})
	})
}
