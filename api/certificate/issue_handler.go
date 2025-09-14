package certificate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/core/certificate"
)

type issueHandler struct {
	commands *certificate.Commands
}

func (h issueHandler) handle(ctx *gin.Context) {
	payload := &issueCertificateRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domainModel := toIssueCertificateRequest(payload)

	cert, err := h.commands.Issue(ctx.Request.Context(), domainModel)
	if api_error.CanHandle(err) {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toIssueCertificateResponse(cert, err))
}
