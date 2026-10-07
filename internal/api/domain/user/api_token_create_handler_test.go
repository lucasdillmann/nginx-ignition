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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/coreerror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

func Test_apiTokenCreateHandler(t *testing.T) {
	setup := func(t *testing.T, usr *user.User) (*user.MockedCommands, *gin.Engine) {
		controller := gomock.NewController(t)
		commands := user.NewMockedCommands(controller)
		authorizer, _ := authorization.New(configuration.New(), commands)
		handler := apiTokenCreateHandler{
			commands:   commands,
			authorizer: authorizer,
		}

		engine := gin.New()
		engine.Use(func(ginContext *gin.Context) {
			ginContext.Set("ABAC:Subject", &authorization.Subject{User: usr})
			ginContext.Next()
		})
		engine.POST("/current/tokens", handler.handle)

		return commands, engine
	}

	performRequest := func(engine *gin.Engine, payload any) *httptest.ResponseRecorder {
		body, _ := json.Marshal(payload)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/current/tokens", bytes.NewBuffer(body))
		engine.ServeHTTP(recorder, request)
		return recorder
	}

	t.Run("handle", func(t *testing.T) {
		t.Run("returns 201 Created with the generated token", func(t *testing.T) {
			usr := newUser()
			token := newAPIToken(usr)
			payload := apiTokenCreateRequestDTO{
				Name:       new("automation"),
				Expiration: token.Expiration,
			}

			commands, engine := setup(t, usr)
			commands.EXPECT().
				CreateAPIToken(gomock.Any(), usr.ID, gomock.Any()).
				DoAndReturn(func(_ any, _ uuid.UUID, request *user.NewAPITokenRequest) (*user.APIToken, error) {
					assert.Equal(t, "automation", request.Name)
					assert.Equal(t, token.Expiration, request.Expiration)
					return token, nil
				})

			recorder := performRequest(engine, payload)

			assert.Equal(t, http.StatusCreated, recorder.Code)
			var response apiTokenCreatedResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, token.ID, response.ID)
			assert.Equal(t, token.Name, response.Name)
			assert.NotEmpty(t, response.Token)
		})

		t.Run("returns 201 Created for a token without expiration", func(t *testing.T) {
			usr := newUser()
			token := newAPIToken(usr)
			token.Expiration = nil
			payload := apiTokenCreateRequestDTO{Name: new("automation")}

			commands, engine := setup(t, usr)
			commands.EXPECT().
				CreateAPIToken(gomock.Any(), usr.ID, gomock.Any()).
				Return(token, nil)

			recorder := performRequest(engine, payload)

			assert.Equal(t, http.StatusCreated, recorder.Code)
			var response apiTokenCreatedResponseDTO
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Nil(t, response.Expiration)
			assert.NotEmpty(t, response.Token)
		})

		t.Run("panics when the token cannot be created", func(t *testing.T) {
			usr := newUser()
			payload := apiTokenCreateRequestDTO{Name: new("")}

			commands, engine := setup(t, usr)
			commands.EXPECT().
				CreateAPIToken(gomock.Any(), usr.ID, gomock.Any()).
				Return(nil, coreerror.New(i18n.M(t.Context(), i18n.K.CoreUserTokenExpiredDate), true))

			assert.Panics(t, func() {
				performRequest(engine, payload)
			})
		})

		t.Run("returns 401 Unauthorized without a subject", func(t *testing.T) {
			controller := gomock.NewController(t)
			commands := user.NewMockedCommands(controller)
			handler := apiTokenCreateHandler{commands: commands}
			engine := gin.New()
			engine.POST("/current/tokens", handler.handle)

			recorder := performRequest(engine, apiTokenCreateRequestDTO{Name: new("automation")})

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("returns 400 Bad Request when called with an API token", func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			usr := newUser()
			commands := user.NewMockedCommands(controller)
			handler := apiTokenCreateHandler{commands: commands}

			engine := gin.New()
			engine.Use(func(ginContext *gin.Context) {
				ginContext.Set("ABAC:Subject", &authorization.Subject{
					User: usr,
					Kind: authorization.APIKind,
				})
				ginContext.Next()
			})
			engine.POST("/current/tokens", handler.handle)

			recorder := performRequest(engine, apiTokenCreateRequestDTO{Name: new("automation")})

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
	})
}
