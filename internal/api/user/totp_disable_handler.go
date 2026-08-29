package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/user"
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
