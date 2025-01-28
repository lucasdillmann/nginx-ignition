package certificate

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type availableProvidersHandler struct {
	command *certificate.AvailableProvidersCommand
}

func (h availableProvidersHandler) handle(context *gin.Context) {
	availableProviders, err := (*h.command)()
	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, toAvailableProviderResponse(availableProviders))
}
