package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/certificate/server"
)

type issueHandler struct {
	commands *server.Commands
}

func (h issueHandler) handle(ctx *gin.Context) {
	payload := &issueCertificateRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(toIssueCertificateRequest, payload)

	cert, err := h.commands.Issue(ctx.Request.Context(), domainModel)
	if apierror.CanHandle(err) {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toIssueCertificateResponse(cert, err))
}
