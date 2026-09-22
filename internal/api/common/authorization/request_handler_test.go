package authorization

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

type requestHandlerSetup struct {
	authorizer *ABAC
	commands   *user.MockedCommands
	engine     *gin.Engine
}

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_RequestHandler(t *testing.T) {
	t.Run("handles", func(t *testing.T) {
		t.Run("passes through for paths outside the API", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})

			recorder := performRequest(setup.engine, "GET", "/public", "")

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("allows anonymous paths", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})

			recorder := performRequest(setup.engine, "GET", "/api/anonymous", "")

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("rejects a missing access token", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", "")

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("rejects an invalid access token", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", "invalid-token")

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("allows paths allowed for all users with a valid token", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			token, _ := setup.authorizer.Jwt().GenerateToken(usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/all-users", *token)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("denies access without the required permission", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			token, _ := setup.authorizer.Jwt().GenerateToken(usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", *token)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("denies write methods with read-only access", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			token, _ := setup.authorizer.Jwt().GenerateToken(usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "POST", "/api/nginx/status", *token)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("grants access and stores the subject in the context", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			token, _ := setup.authorizer.Jwt().GenerateToken(usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", *token)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, usr.ID.String(), recorder.Header().Get("X-Subject"))
		})

		t.Run("refreshes the token when within the renewal window", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{
				"nginx-ignition.security.jwt.renew-window-seconds": "7200",
			})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			token, _ := setup.authorizer.Jwt().GenerateToken(usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", *token)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.NotEmpty(t, recorder.Header().Get("Authorization"))
			assert.NotEqual(t, "Bearer "+*token, recorder.Header().Get("Authorization"))
		})
	})
}

func newRequestHandlerSetup(t *testing.T, overrides map[string]string) *requestHandlerSetup {
	t.Helper()
	authorizer, commands := newAuthorizerWithOverrides(t, overrides)

	engine := gin.New()
	engine.Use(gin.CustomRecoveryWithWriter(nil, apierror.Handler))
	engine.Use(authorizer.HandleRequest)

	protectedGroup := authorizer.ConfigureGroup(
		engine,
		"/api/nginx",
		func(permissions user.Permissions) user.AccessLevel { return permissions.NginxServer },
	)
	protectedGroup.Any("/status", func(ctx *gin.Context) {
		subject := CurrentSubject(ctx)
		ctx.Header("X-Subject", subject.User.ID.String())
		ctx.Status(http.StatusOK)
	})

	authorizer.AllowAnonymous("GET", "/api/anonymous")
	engine.GET("/api/anonymous", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	authorizer.AllowAllUsers("GET", "/api/all-users")
	engine.GET("/api/all-users", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	engine.GET("/public", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	return &requestHandlerSetup{
		authorizer: authorizer,
		commands:   commands,
		engine:     engine,
	}
}

func performRequest(engine *gin.Engine, method, path, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	return recorder
}
