package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

type totpActivateHandler struct {
	commands user.Commands
}

func (h totpActivateHandler) handle(ctx *gin.Context) {
	payload := &totpActivateRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	currentUserID := authorization.CurrentSubject(ctx).User.ID

	activated, err := h.commands.ActivateTOTP(
		ctx.Request.Context(),
		currentUserID,
		*payload.Code,
	)
	if err != nil {
		panic(err)
	}

	if !activated {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}
