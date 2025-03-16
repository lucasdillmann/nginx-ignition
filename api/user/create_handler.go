package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

type createHandler struct {
	command *user.SaveCommand
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &userRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	domainModel := toDomain(payload)
	domainModel.ID = uuid.New()
	currentUserId := authorization.CurrentSubject(ctx).User.ID

	if err := (*h.command)(ctx.Request.Context(), domainModel, &currentUserId); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
