package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

type totpDisableHandler struct {
	commands user.Commands
}

func (h totpDisableHandler) handle(ctx *gin.Context) {
	currentUserID := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.DisableTOTP(ctx.Request.Context(), currentUserID); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
