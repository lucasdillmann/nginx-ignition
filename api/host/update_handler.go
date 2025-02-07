package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

type updateHandler struct {
	command *host.SaveCommand
}

func (h updateHandler) handle(context *gin.Context) {
	payload := &hostRequestDto{}
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

	if err = (*h.command)(domainModel); err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
