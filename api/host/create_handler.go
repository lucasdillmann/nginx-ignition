package host

import (
	"dillmann.com.br/nginx-ignition/core/host"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

type createHandler struct {
	commands *host.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &hostRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domainModel := toDomain(payload)
	domainModel.ID = uuid.New()

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
