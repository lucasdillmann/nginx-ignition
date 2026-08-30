package nginx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/core/nginx"
	"github.com/lucasdillmann/nginx-ignition/internal/core/settings"
)

type metadataHandler struct {
	nginxCommands    nginx.Commands
	settingsCommands settings.Commands
}

func (h metadataHandler) handle(ctx *gin.Context) {
	metadata, err := h.nginxCommands.GetMetadata(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	set, err := h.settingsCommands.Get(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"version": metadata.Version,
		"stats": gin.H{
			"enabled":  set.Nginx.Stats.Enabled,
			"allHosts": set.Nginx.Stats.AllHosts,
		},
		"availableSupport": gin.H{
			"streams": metadata.StreamSupportType(),
			"runCode": metadata.RunCodeSupportType(),
			"tlsSni":  metadata.SNISupportType(),
			"stats":   metadata.StatsSupportType(),
		},
	})
}
