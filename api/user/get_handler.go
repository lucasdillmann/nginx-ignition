package user

import (
	"net/http"

	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getHandler struct {
	commands *user.Commands
}

func (h getHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	usr, err := h.commands.Get(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if usr == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, toDto(usr))
}
