package user_api

import (
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type onboardingFinishHandler struct {
	statusCommand       *user.OnboardingCompletedCommand
	saveCommand         *user.SaveCommand
	authenticateCommand *user.AuthenticateCommand
	authorizer          *authorization.RBAC
}

func (h onboardingFinishHandler) handle(context *gin.Context) {
	alreadyFinished, err := (*h.statusCommand)()
	if err != nil {
		panic(err)
	}

	if alreadyFinished {
		context.Status(http.StatusForbidden)
		return
	}

	requestPayload := &userRequestDto{}
	if err = context.BindJSON(requestPayload); err != nil {
		panic(err)
	}

	domainModel := toDomain(requestPayload)
	domainModel.ID = uuid.New()
	domainModel.Enabled = true
	domainModel.Role = user.AdminRole

	if err = (*h.saveCommand)(domainModel, nil); err != nil {
		panic(err)
	}

	usr, err := (*h.authenticateCommand)(domainModel.Username, *domainModel.Password)
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
