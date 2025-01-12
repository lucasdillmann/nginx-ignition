package user_api

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
)

type logoutHandler struct {
	authorizer *authorization.RBAC
}

func (h logoutHandler) handle(context *gin.Context) {
	subject := authorization.CurrentSubject(context)
	h.authorizer.Jwt().RevokeToken(subject.TokenID)
	context.Status(http.StatusNoContent)
}
