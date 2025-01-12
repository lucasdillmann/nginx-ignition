package user_api

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type getHandler struct {
	command *user.GetCommand
}

func (h getHandler) handle(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	usr, err := (*h.command)(id)
	if err != nil {
		panic(err)
	}

	if usr == nil {
		context.Status(http.StatusNotFound)
		return
	}

	context.JSON(http.StatusOK, toDto(usr))
}
