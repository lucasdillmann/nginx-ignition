package certificate

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type issueHandler struct {
	command *certificate.IssueCommand
}

func (h issueHandler) handle(context *gin.Context) {
	payload := &issueCertificateRequest{}
	if err := context.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domainModel := toIssueCertificateRequest(payload)

	cert, err := (*h.command)(domainModel)
	if api_error.CanHandle(err) {
		panic(err)
	}

	context.JSON(http.StatusOK, toIssueCertificateResponse(cert, err))
}
