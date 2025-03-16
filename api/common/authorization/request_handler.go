package authorization

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	RequestSubject = "RBAC:Subject"
)

func (m *RBAC) HandleRequest(ctx *gin.Context) {
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

	requiredRole := m.findRequiredRole(ctx.Request.Method, path)
	if requiredRole != nil && subject.User.Role != *requiredRole {
		ctx.Abort()
		panic(api_error.New(
			http.StatusForbidden,
			"User does not have the required role to access this resource",
		))
	}

	refreshedToken, _ := m.jwt.RefreshToken(subject)
	if refreshedToken != nil {
		ctx.Header("Authorization", "Bearer "+*refreshedToken)
	}

	ctx.Set(RequestSubject, subject)
	ctx.Next()
}
