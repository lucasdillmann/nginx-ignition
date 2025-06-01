package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
)

type logoutHandler struct {
	authorizer *authorization.ABAC
}

func (h logoutHandler) handle(ctx *gin.Context) {
	subject := authorization.CurrentSubject(ctx)
	h.authorizer.Jwt().RevokeToken(subject.TokenID)
	ctx.Status(http.StatusNoContent)
}
