package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
)

type loginHandler struct {
	commands   user.Commands
	authorizer *authorization.ABAC
}

func (h loginHandler) handle(ctx *gin.Context) {
	requestPayload := &userLoginRequestDTO{}
	if err := ctx.BindJSON(&requestPayload); err != nil {
		panic(err)
	}

	totp := requestPayload.TOTP
	if totp == nil {
		totp = new("")
	}

	outcome, usr, err := h.commands.Authenticate(
		ctx.Request.Context(),
		*requestPayload.Username,
		*requestPayload.Password,
		*totp,
	)
	if err != nil {
		panic(err)
	}

	if outcome != user.AuthenticationSuccessful || usr == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"reason": outcome,
		})

		return
	}

	token, err := h.authorizer.Jwt().GenerateToken(usr)
	if err != nil {
		panic(err)
	}

	responsePayload := &userLoginResponseDTO{*token}
	ctx.JSON(http.StatusOK, responsePayload)
}
