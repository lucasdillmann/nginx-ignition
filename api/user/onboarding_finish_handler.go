package user

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

func (h onboardingFinishHandler) handle(ctx *gin.Context) {
	alreadyFinished, err := (*h.statusCommand)(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	if alreadyFinished {
		ctx.Status(http.StatusForbidden)
		return
	}

	requestPayload := &userRequestDto{}
	if err = ctx.BindJSON(requestPayload); err != nil {
		panic(err)
	}

	domainModel := toDomain(requestPayload)
	domainModel.ID = uuid.New()
	domainModel.Enabled = true
	domainModel.Role = user.AdminRole

	if err = (*h.saveCommand)(ctx.Request.Context(), domainModel, nil); err != nil {
		panic(err)
	}

	usr, err := (*h.authenticateCommand)(ctx.Request.Context(), domainModel.Username, *domainModel.Password)
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
