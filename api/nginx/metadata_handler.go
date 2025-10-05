package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/nginx"
)

type metadataHandler struct {
	commands *nginx.Commands
}

func (h metadataHandler) handle(ctx *gin.Context) {
	metadata, err := h.commands.GetMetadata(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"version":      metadata.Version,
		"buildDetails": metadata.BuildDetails,
		"modules":      metadata.Modules,
		"availableSupport": gin.H{
			"streams": metadata.StreamSupportType(),
			"runCode": metadata.RunCodeSupportAvailable(),
			"tlsSni":  metadata.TLSSNIEnabled,
		},
	})
}
