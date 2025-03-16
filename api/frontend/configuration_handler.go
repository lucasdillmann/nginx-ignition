package frontend

import (
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"github.com/gin-gonic/gin"
	"net/http"
)

type configurationHandler struct {
	configuration *configuration.Configuration
}

func (h *configurationHandler) handle(ctx *gin.Context) {
	codeEditorApiKey, _ := h.configuration.Get("nginx-ignition.frontend.code-editor-api-key")

	var apiKey *string
	if codeEditorApiKey != "" {
		apiKey = &codeEditorApiKey
	}

	output := &configurationDto{
		CodeEditor: codeEditorDto{
			ApiKey: apiKey,
		},
	}

	ctx.JSON(http.StatusOK, output)
}
