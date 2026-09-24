package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type loginHandler struct {
	commands   user2.Commands
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

	if requestPayload.Username == nil || requestPayload.Password == nil {
		ctx.Status(http.StatusBadRequest)
		return
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

	if outcome != user2.AuthenticationSuccessful || usr == nil {
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
