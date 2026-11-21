package certificate

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/converter"
	"github.com/lucasdillmann/nginx-ignition/internal/api/core/eventstream"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

type issueHandler struct {
	commands certificate.Commands
}

func (h issueHandler) handle(ctx *gin.Context) {
	payload := &issueCertificateRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(ctx.Request.Context(), toIssueCertificateRequest, payload)

	eventstream.Serve(ctx, func(runContext context.Context) (any, error) {
		cert, err := h.commands.Issue(runContext, domainModel)

		return toIssueCertificateResponse(cert, err), err
	})
}
