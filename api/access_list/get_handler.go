package access_list

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/access_list"
)

type getHandler struct {
	commands *access_list.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	accessList, err := h.commands.Get(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if accessList == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toDto(accessList))
}
