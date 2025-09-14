package host

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/host"
)

type updateHandler struct {
	commands *host.Commands
}

func (h updateHandler) handle(ctx *gin.Context) {
	payload := &hostRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil || id == uuid.Nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	domainModel := toDomain(payload)
	domainModel.ID = id

	if err = h.commands.Save(ctx.Request.Context(), domainModel); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
