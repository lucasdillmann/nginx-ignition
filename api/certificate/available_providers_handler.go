package certificate

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type availableProvidersHandler struct {
	command *certificate.AvailableProvidersCommand
}

func (h availableProvidersHandler) handle(ctx *gin.Context) {
	availableProviders, err := (*h.command)(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toAvailableProviderResponse(availableProviders))
}
