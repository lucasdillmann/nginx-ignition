package certificate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/api/common/converter"
	"github.com/lucasdillmann/nginx-ignition/internal/core/certificate"
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

	cert, err := h.commands.Issue(ctx.Request.Context(), domainModel)
	if apierror.CanHandle(err) {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toIssueCertificateResponse(cert, err))
}
