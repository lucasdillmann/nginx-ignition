package user_api

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

type updateHandler struct {
	command *user.SaveCommand
}

func (h updateHandler) handle(context *gin.Context) {
	payload := &userRequestDto{}
	if err := context.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	id, err := uuid.Parse(context.Param("id"))
	if err != nil || id == uuid.Nil {
		context.Status(http.StatusNotFound)
		return
	}

	domainModel := toDomain(payload)
	domainModel.ID = id
	currentUserId := authorization.CurrentSubject(context).User.ID

	if err := (*h.command)(domainModel, &currentUserId); err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
