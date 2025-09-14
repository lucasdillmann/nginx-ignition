package certificate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

type availableProvidersHandler struct {
	commands *certificate.Commands
}

func (h availableProvidersHandler) handle(ctx *gin.Context) {
	availableProviders, err := h.commands.AvailableProviders(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, toAvailableProviderResponse(availableProviders))
}
