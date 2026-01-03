package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

type updatePasswordHandler struct {
	commands user.Commands
}

func (h updatePasswordHandler) handle(ctx *gin.Context) {
	payload := &userPasswordUpdateRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	currentUserID := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.UpdatePassword(ctx.Request.Context(), currentUserID, *payload.CurrentPassword, *payload.NewPassword); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
