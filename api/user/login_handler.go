package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type loginHandler struct {
	commands   *user.Commands
	authorizer *authorization.ABAC
}

func (h loginHandler) handle(ctx *gin.Context) {
	requestPayload := &userLoginRequestDto{}
	if err := ctx.BindJSON(&requestPayload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(requestPayload); err != nil {
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

	responsePayload := &userLoginResponseDto{*token}
	ctx.JSON(http.StatusOK, responsePayload)
}
