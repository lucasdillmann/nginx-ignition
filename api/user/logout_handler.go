package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
)

type logoutHandler struct {
	authorizer *authorization.ABAC
}

func (h logoutHandler) handle(ctx *gin.Context) {
	subject := authorization.CurrentSubject(ctx)
	h.authorizer.Jwt().RevokeToken(subject.TokenID)
	ctx.Status(http.StatusNoContent)
}
