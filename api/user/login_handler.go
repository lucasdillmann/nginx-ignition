package user

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type loginHandler struct {
	command    *user.AuthenticateCommand
	authorizer *authorization.RBAC
}

func (h loginHandler) handle(context *gin.Context) {
	requestPayload := &userLoginRequestDto{}
	if err := context.BindJSON(&requestPayload); err != nil {
		panic(err)
	}

	if err := validator.New().Struct(requestPayload); err != nil {
		panic(err)
	}

	usr, err := (*h.command)(*requestPayload.Username, *requestPayload.Password)
	if err != nil {
		panic(err)
	}

	token, err := h.authorizer.Jwt().GenerateToken(usr)
	if err != nil {
		panic(err)
	}

	responsePayload := &userLoginResponseDto{*token}
	context.JSON(http.StatusOK, responsePayload)
}
