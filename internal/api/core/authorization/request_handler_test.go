package authorization

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
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
			token := setup.sessionToken(t, usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/all-users", token)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run("denies access without the required permission", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			token := setup.sessionToken(t, usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", token)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("denies write methods with read-only access", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			token := setup.sessionToken(t, usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "POST", "/api/nginx/status", token)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("grants access and stores the subject in the context", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			token := setup.sessionToken(t, usr)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", token)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, usr.ID.String(), recorder.Header().Get("X-Subject"))
		})

		t.Run("refreshes the token when within the renewal window", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{
				"nginx-ignition.security.jwt.renew-window-seconds": "30",
			})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			session := newSession(usr)
			token := signRawToken(t, setup.authorizer, jwt.MapClaims{
				"aud":          uniqueIdentifier,
				"iss":          uniqueIdentifier,
				"nbf":          time.Now().Add(time.Second * -60).Unix(),
				"iat":          time.Now().Add(time.Second * -60).Unix(),
				"jti":          session.IterationID.String(),
				sessionIDClaim: session.ID.String(),
				"sub":          usr.ID.String(),
				tokenKindClaim: SessionKind,
				"exp":          time.Now().Add(time.Second * 5).Unix(),
			})

			rotated := newSession(usr)

			setup.commands.EXPECT().FindSession(gomock.Any(), session.ID).Return(session, nil)
			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
			setup.commands.EXPECT().
				RotateSession(gomock.Any(), session.ID, usr.ID, session.IterationID, gomock.Any()).
				Return(rotated, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", token)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.NotEmpty(t, recorder.Header().Get("Authorization"))
			assert.NotEqual(t, "Bearer "+token, recorder.Header().Get("Authorization"))
			assert.Equal(t, session.ID.String(), recorder.Header().Get("X-Session"))
		})

		t.Run("never refreshes an API token", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{
				"nginx-ignition.security.jwt.renew-window-seconds": "30",
			})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			tokenID := uuid.New()
			token, _ := setup.authorizer.Jwt().GenerateToken(
				context.Background(),
				usr,
				APIKind,
				&tokenID,
				nil,
			)

			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
			setup.commands.EXPECT().
				FindAPIToken(gomock.Any(), usr.ID, tokenID).
				Return(&user.APIToken{ID: tokenID, UserID: usr.ID}, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", *token)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Empty(t, recorder.Header().Get("Authorization"))
		})

		t.Run("rejects an API token that was revoked", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
			tokenID := uuid.New()
			token, _ := setup.authorizer.Jwt().GenerateToken(
				context.Background(),
				usr,
				APIKind,
				&tokenID,
				nil,
			)

			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
			setup.commands.EXPECT().FindAPIToken(gomock.Any(), usr.ID, tokenID).Return(nil, nil)

			recorder := performRequest(setup.engine, "GET", "/api/nginx/status", *token)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
		})

		t.Run("rejects an API token on a sessions-tokens-only path", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			tokenID := uuid.New()
			token, _ := setup.authorizer.Jwt().GenerateToken(
				context.Background(),
				usr,
				APIKind,
				&tokenID,
				nil,
			)

			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)
			setup.commands.EXPECT().
				FindAPIToken(gomock.Any(), usr.ID, tokenID).
				Return(&user.APIToken{ID: tokenID, UserID: usr.ID}, nil)

			recorder := performRequest(setup.engine, "POST", "/api/session-only", *token)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
		})

		t.Run("allows a session token on a sessions-tokens-only path", func(t *testing.T) {
			setup := newRequestHandlerSetup(t, map[string]string{})
			usr := newUser()
			token := setup.sessionToken(t, usr)

			setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

			recorder := performRequest(setup.engine, "POST", "/api/session-only", token)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})

		t.Run(
			"still applies the permission check on sessions-tokens-only paths",
			func(t *testing.T) {
				setup := newRequestHandlerSetup(t, map[string]string{})
				usr := newUser()
				usr.Permissions.NginxServer = user.NoAccessAccessLevel
				token := setup.sessionToken(t, usr)

				setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

				recorder := performRequest(setup.engine, "GET", "/api/restricted/status", token)

				assert.Equal(t, http.StatusForbidden, recorder.Code)
			},
		)

		t.Run(
			"allows a session token with permission on sessions-tokens-only paths",
			func(t *testing.T) {
				setup := newRequestHandlerSetup(t, map[string]string{})
				usr := newUser()
				usr.Permissions.NginxServer = user.ReadOnlyAccessLevel
				token := setup.sessionToken(t, usr)

				setup.commands.EXPECT().Get(gomock.Any(), usr.ID).Return(usr, nil)

				recorder := performRequest(setup.engine, "GET", "/api/restricted/status", token)

				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		)
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
		ctx.Header("X-Session", subject.SessionID.String())
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

	authorizer.AllowAllUsers("POST", "/api/session-only")
	authorizer.AllowSessionsTokensOnly("POST", "/api/session-only")
	engine.POST("/api/session-only", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	restrictedGroup := authorizer.ConfigureGroup(
		engine,
		"/api/restricted",
		func(permissions user.Permissions) user.AccessLevel { return permissions.NginxServer },
	)
	authorizer.AllowSessionsTokensOnly("GET", "/api/restricted/status")
	restrictedGroup.GET("/status", func(ctx *gin.Context) {
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

func (s *requestHandlerSetup) sessionToken(t *testing.T, usr *user.User) string {
	t.Helper()

	session := newSession(usr)
	token := issueSessionToken(t, s.authorizer, s.commands, usr, session)

	s.commands.EXPECT().FindSession(gomock.Any(), session.ID).Return(session, nil)

	return *token
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
