package vpn

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/converter"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/vpn"
)

type putHandler struct {
	commands vpn.Commands
}

func (h putHandler) handle(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Status(http.StatusNotFound)
		return
	}

	uuidValue, err := uuid.Parse(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	payload := &vpnRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	data := converter.Wrap2(ctx.Request.Context(), fromDTO, uuidValue, payload)

	if err := h.commands.Save(ctx.Request.Context(), data); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
