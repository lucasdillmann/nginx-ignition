package backup

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/backup"
)

type getHandler struct {
	commands *backup.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	data, err := h.commands.Get(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", data.FileName))
	ctx.Data(http.StatusOK, data.ContentType, data.Contents)
}
