package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
)

type currentHandler struct{}

func (h currentHandler) handle(ctx *gin.Context) {
	currentUser := authorization.CurrentSubject(ctx).User
	ctx.JSON(http.StatusOK, toDto(currentUser))
}
