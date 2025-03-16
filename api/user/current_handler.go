package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
)

type currentHandler struct {
}

func (h currentHandler) handle(ctx *gin.Context) {
	currentUser := authorization.CurrentSubject(ctx).User
	ctx.JSON(http.StatusOK, toDto(currentUser))
}
