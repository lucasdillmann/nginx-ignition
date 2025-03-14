package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type deleteHandler struct {
	command *access_list.DeleteCommand
}

func (h deleteHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	err = (*h.command)(id)
	if err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
