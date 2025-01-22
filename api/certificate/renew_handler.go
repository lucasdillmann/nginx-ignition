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

func (h renewHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	err = (*h.command)(id)
	if api_error.CanHandle(err) {
		panic(err)
	}

	context.JSON(http.StatusOK, toRenewCertificateResponse(err))
}
