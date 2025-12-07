package vpn

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type createHandler struct {
	commands *vpn.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &vpnRequest{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap2(toDomain, payload, uuid.New())
	if err := h.commands.Save(ctx.Request.Context(), domainModel); err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusCreated,
		map[string]any{
			"id": domainModel.ID,
		},
	)
}
