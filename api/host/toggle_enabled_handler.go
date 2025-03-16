package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type toggleEnabledHandler struct {
	getCommand  *host.GetCommand
	saveCommand *host.SaveCommand
}

func (h toggleEnabledHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	data, err := (*h.getCommand)(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if data == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	data.Enabled = !data.Enabled
	err = (*h.saveCommand)(ctx.Request.Context(), data)

	if err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
