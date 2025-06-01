package authorization

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
		panic(errInvalidToken)
	}

	subject, err := m.jwt.ValidateToken(ctx.Request.Context(), accessToken)
	if err != nil {
		ctx.Abort()
		panic(api_error.New(
			http.StatusUnauthorized,
			"Invalid or expired access token",
		))
	}

	if !m.isAllowedForAllUsers(ctx.Request.Method, path) {
		accessGranted := m.isAccessGranted(ctx.Request.Method, path, &subject.User.Permissions)
		if !accessGranted {
			ctx.Abort()
			panic(api_error.New(
				http.StatusForbidden,
				"User does not have the required permission to access this resource",
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
