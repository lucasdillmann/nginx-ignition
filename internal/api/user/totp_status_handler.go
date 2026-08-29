package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

type totpStatusHandler struct {
	commands user.Commands
}

func (h totpStatusHandler) handle(ctx *gin.Context) {
	currentUserID := authorization.CurrentSubject(ctx).User.ID

	enabled, err := h.commands.GetTOTPStatus(ctx.Request.Context(), currentUserID)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, &totpStatusResponseDTO{enabled})
}
