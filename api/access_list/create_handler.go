package access_list

import (
	"dillmann.com.br/nginx-ignition/core/access_list"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

type createHandler struct {
	command *access_list.SaveCommand
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &accessListRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domainModel := toDomain(payload)
	domainModel.ID = uuid.New()

	if err := (*h.command)(ctx.Request.Context(), domainModel); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
