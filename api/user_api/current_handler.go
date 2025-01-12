package user_api

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"github.com/gin-gonic/gin"
	"net/http"
)

type currentHandler struct {
}

func (h currentHandler) handle(context *gin.Context) {
	currentUser := authorization.CurrentSubject(context).User
	context.JSON(http.StatusOK, toDto(currentUser))
}
