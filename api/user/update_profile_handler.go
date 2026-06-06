package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	coreuser "dillmann.com.br/nginx-ignition/core/user"
)

type updateProfileHandler struct {
	commands coreuser.Commands
}

func (h updateProfileHandler) handle(ctx *gin.Context) {
	payload := &userProfileUpdateRequestDTO{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	currentUserID := authorization.CurrentSubject(ctx).User.ID

	if err := h.commands.UpdateProfile(
		ctx.Request.Context(),
		currentUserID,
		getStringValue(payload.Name),
		getStringValue(payload.Username),
		payload.NotificationLanguage,
	); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
