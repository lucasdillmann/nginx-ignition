package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
)

type logoutHandler struct {
	authorizer *authorization.RBAC
}

func (h logoutHandler) handle(ctx *gin.Context) {
	subject := authorization.CurrentSubject(ctx)
	h.authorizer.Jwt().RevokeToken(subject.TokenID)
	ctx.Status(http.StatusNoContent)
}
