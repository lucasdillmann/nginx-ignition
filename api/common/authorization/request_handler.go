package authorization

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

const (
	RequestSubject = "ABAC:Subject"
)

func (m *ABAC) HandleRequest(ctx *gin.Context) {
	path := ctx.FullPath()
	if !strings.HasPrefix(path, "/api/") {
		ctx.Next()
		return
	}

	if m.isAnonymous(ctx.Request.Method, path) {
		ctx.Next()
		return
	}

	accessToken, _ := strings.CutPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	if strings.TrimSpace(accessToken) == "" {
		ctx.Abort()
		panic(apierror.New(
			http.StatusUnauthorized,
			i18n.M(ctx.Request.Context(), i18n.K.AuthorizationErrorInvalidAccessToken),
		))
	}

	subject, err := m.jwt.ValidateToken(ctx.Request.Context(), accessToken)
	if err != nil {
		ctx.Abort()
		panic(apierror.New(
			http.StatusUnauthorized,
			i18n.M(ctx.Request.Context(), i18n.K.AuthorizationErrorInvalidAccessToken),
		))
	}

	if !m.isAllowedForAllUsers(ctx.Request.Method, path) {
		accessGranted := m.isAccessGranted(ctx.Request.Method, path, &subject.User.Permissions)
		if !accessGranted {
			ctx.Abort()
			panic(apierror.New(
				http.StatusForbidden,
				i18n.M(ctx.Request.Context(), i18n.K.AuthorizationErrorAccessDenied),
			))
		}
	}

	refreshedToken, _ := m.jwt.RefreshToken(subject)
	if refreshedToken != nil {
		ctx.Header("Authorization", "Bearer "+*refreshedToken)
	}

	ctx.Set(RequestSubject, subject)
	ctx.Next()
}
