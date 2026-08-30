package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
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

	if err := h.commands.UpdatePassword(
		ctx.Request.Context(),
		currentUserID,
		*payload.CurrentPassword,
		*payload.NewPassword,
	); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
