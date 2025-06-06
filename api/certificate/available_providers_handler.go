package certificate

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"github.com/gin-gonic/gin"
	"net/http"
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
