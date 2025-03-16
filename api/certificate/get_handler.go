package certificate

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type getHandler struct {
	command *certificate.GetCommand
}

func (h getHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	accessList, err := (*h.command)(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if accessList == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toCertificateResponse(accessList))
}
