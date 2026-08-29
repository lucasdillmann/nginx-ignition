package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

type totpEnableHandler struct {
	commands user.Commands
}

func (h totpEnableHandler) handle(ctx *gin.Context) {
	currentUserID := authorization.CurrentSubject(ctx).User.ID

	alreadyActivated, err := h.commands.GetTOTPStatus(ctx.Request.Context(), currentUserID)
	if err != nil {
		panic(err)
	}

	if alreadyActivated {
		ctx.Status(http.StatusBadRequest)
		return
	}

	secret, err := h.commands.EnableTOTP(ctx.Request.Context(), currentUserID)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, &totpEnableResponseDTO{secret})
}
