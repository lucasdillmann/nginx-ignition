package certificate

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type renewHandler struct {
	command *certificate.RenewCommand
}

func (h renewHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	err = (*h.command)(ctx.Request.Context(), id)
	if api_error.CanHandle(err) {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toRenewCertificateResponse(err))
}
