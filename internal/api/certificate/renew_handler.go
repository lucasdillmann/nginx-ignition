package certificate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
)

type renewHandler struct {
	commands certificate.Commands
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
