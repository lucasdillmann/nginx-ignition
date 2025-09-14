package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

type deleteHandler struct {
	commands *user.Commands
}

func (h deleteHandler) handle(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	currentUserId := authorization.CurrentSubject(ctx).User.ID
	if id == currentUserId {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "You cannot delete your own user",
			},
		)
		return
	}

	err = h.commands.Delete(ctx.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
