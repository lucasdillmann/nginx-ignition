package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
)

type currentHandler struct{}

func (h currentHandler) handle(ctx *gin.Context) {
	currentSubject := authorization.CurrentSubject(ctx)
	if currentSubject == nil || currentSubject.User == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, toDTO(currentSubject.User))
}
