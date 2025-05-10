package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)

type updateHandler struct {
	commands *user.Commands
}

func (h updateHandler) handle(ctx *gin.Context) {
	payload := &userRequestDto{}
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
	currentUserId := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.Save(ctx.Request.Context(), domainModel, &currentUserId); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
