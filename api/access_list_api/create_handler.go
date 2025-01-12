package access_list_api

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

func (h createHandler) handle(context *gin.Context) {
	payload := &accessListRequestDto{}
	if err := context.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domainModel := toDomain(payload)
	domainModel.ID = uuid.New()

	if err := (*h.command)(domainModel); err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
