package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/user"
)

type createHandler struct {
	commands *user.Commands
}

func (h createHandler) handle(ctx *gin.Context) {
	payload := &userRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(toDomain, payload)
	domainModel.ID = uuid.New()
	currentUserID := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.Save(ctx.Request.Context(), domainModel, &currentUserID); err != nil {
		panic(err)
	}

	ctx.JSON(
		http.StatusCreated,
		map[string]any{
			"id": domainModel.ID,
		},
	)
}
