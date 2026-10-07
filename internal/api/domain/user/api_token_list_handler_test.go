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
	"github.com/lucasdillmann/nginx-ignition/internal/api/core/pagination"

	corepagination "github.com/lucasdillmann/nginx-ignition/internal/business/core/pagination"
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
		t.Run("returns 200 OK with the requested page of user tokens", func(t *testing.T) {
			usr := newUser()
			token := newAPIToken(usr)

			commands, engine := setup(t, usr.ID)
			commands.EXPECT().
				ListAPITokens(gomock.Any(), usr.ID, 25, 0, nil).
				Return(corepagination.New(0, 25, 1, []user.APIToken{*token}), nil)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/tokens", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response pagination.PageDTO[apiTokenResponseDTO]
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, 1, response.TotalItems)
			require.Len(t, response.Contents, 1)
			assert.Equal(t, token.ID, response.Contents[0].ID)
			assert.Equal(t, token.Name, response.Contents[0].Name)
			assert.Equal(t, token.Expiration, response.Contents[0].Expiration)
		})

		t.Run("forwards pagination and search parameters", func(t *testing.T) {
			usr := newUser()
			searchTerms := "automation"

			commands, engine := setup(t, usr.ID)
			commands.EXPECT().
				ListAPITokens(gomock.Any(), usr.ID, 10, 2, &searchTerms).
				Return(corepagination.New(2, 10, 0, []user.APIToken{}), nil)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				"/current/tokens?pageSize=10&pageNumber=2&searchTerms=automation",
				nil,
			)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("returns 401 Unauthorized without a subject", func(t *testing.T) {
			controller := gomock.NewController(t)
			commands := user.NewMockedCommands(controller)
			handler := apiTokenListHandler{commands: commands}
			engine := gin.New()
			engine.GET("/current/tokens", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/tokens", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("returns 400 Bad Request when called with an API token", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			usr := newUser()
			commands := user.NewMockedCommands(controller)
			handler := apiTokenListHandler{commands: commands}

			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{
					User: usr,
					Kind: authorization.APIKind,
				})
				ginContext.Next()
			})
			engine.GET("/current/tokens", handler.handle)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/current/tokens", nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	})
}
