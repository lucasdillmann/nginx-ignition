package user_api

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type deleteHandler struct {
	command *user.DeleteCommand
}

func (h deleteHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	currentUserId := authorization.CurrentSubject(context).User.ID
	if id == currentUserId {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "You cannot delete your own user",
			},
		)
		return
	}

	err = (*h.command)(id)
	if err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
