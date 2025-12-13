package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
	"dillmann.com.br/nginx-ignition/core/certificate/server"
)

type renewHandler struct {
	commands *server.Commands
}

func (h renewHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	err = h.commands.Renew(ctx.Request.Context(), id)
	if apierror.CanHandle(err) {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toRenewCertificateResponse(err))
}
