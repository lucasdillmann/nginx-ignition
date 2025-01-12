package user_api

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

func (h updatePasswordHandler) handle(context *gin.Context) {
	payload := &userPasswordUpdateRequestDto{}
	if err := context.BindJSON(payload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(payload); err != nil {
		panic(err)
	}

	currentUserId := authorization.CurrentSubject(context).User.ID

	if err := (*h.command)(currentUserId, *payload.CurrentPassword, *payload.NewPassword); err != nil {
		panic(err)
	}

	context.Status(http.StatusNoContent)
}
