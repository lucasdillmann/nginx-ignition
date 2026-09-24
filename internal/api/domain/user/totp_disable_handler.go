package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
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
