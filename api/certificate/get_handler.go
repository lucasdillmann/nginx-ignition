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

func (h getHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	accessList, err := (*h.command)(id)
	if err != nil {
		panic(err)
	}

	if accessList == nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, toCertificateResponse(accessList))
}
