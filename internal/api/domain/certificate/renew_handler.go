package certificate

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/eventstream"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
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

	eventstream.Serve(ctx, func(runContext context.Context) (any, error) {
		err := h.commands.Renew(runContext, id)

		return toRenewCertificateResponse(err), err
	})
}
