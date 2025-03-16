package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type updatePasswordHandler struct {
	command *user.UpdatePasswordCommand
}

func (h updatePasswordHandler) handle(ctx *gin.Context) {
	payload := &userPasswordUpdateRequestDto{}
	if err := ctx.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	currentUserId := authorization.CurrentSubject(ctx).User.ID

	if err := (*h.command)(ctx.Request.Context(), currentUserId, *payload.CurrentPassword, *payload.NewPassword); err != nil {
		panic(err)
	}

	ctx.Status(http.StatusNoContent)
}
