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

func (m *RBAC) HandleRequest(context *gin.Context) {
	path := context.FullPath()
	if !strings.HasPrefix(path, "/api/") {
		context.Next()
		return
	}

	if m.isAnonymous(context.Request.Method, path) {
		context.Next()
		return
	}

	accessToken, _ := strings.CutPrefix(context.GetHeader("Authorization"), "Bearer ")
	if strings.TrimSpace(accessToken) == "" {
		context.Abort()
		panic(errInvalidToken)
	}

	subject, err := m.jwt.ValidateToken(accessToken)
	if err != nil {
		context.Abort()
		panic(api_error.New(
			http.StatusUnauthorized,
			"Invalid or expired access token",
		))
	}

	requiredRole := m.findRequiredRole(context.Request.Method, path)
	if requiredRole != nil && subject.User.Role != *requiredRole {
		context.Abort()
		panic(api_error.New(
			http.StatusForbidden,
			"User does not have the required role to access this resource",
		))
	}

	refreshedToken, _ := m.jwt.RefreshToken(subject)
	if refreshedToken != nil {
		context.Header("Authorization", "Bearer "+*refreshedToken)
	}

	context.Set(RequestSubject, subject)
	context.Next()
}
