package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

type loginHandler struct {
	commands   *user.Commands
	authorizer *authorization.ABAC
}

func (h loginHandler) handle(ctx *gin.Context) {
	requestPayload := &userLoginRequestDTO{}
	if err := ctx.BindJSON(&requestPayload); err != nil {
		panic(err)
	}

	usr, err := h.commands.Authenticate(ctx.Request.Context(), *requestPayload.Username, *requestPayload.Password)
	if err != nil {
		panic(err)
	}

	token, err := h.authorizer.Jwt().GenerateToken(usr)
	if err != nil {
		panic(err)
	}

	responsePayload := &userLoginResponseDTO{*token}
	ctx.JSON(http.StatusOK, responsePayload)
}
