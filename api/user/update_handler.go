package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/user"
)

type updateHandler struct {
	commands user.Commands
}

func (h updateHandler) handle(ctx *gin.Context) {
	payload := &userRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil || id == uuid.Nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	domainModel := converter.Wrap(ctx.Request.Context(), toDomain, payload)
	domainModel.ID = id
	currentUserID := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.Save(ctx.Request.Context(), domainModel, &currentUserID); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
