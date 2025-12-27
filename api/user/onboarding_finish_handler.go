package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/api/common/converter"
	"dillmann.com.br/nginx-ignition/core/user"
)

type onboardingFinishHandler struct {
	commands   *user.Commands
	authorizer *authorization.ABAC
}

func (h onboardingFinishHandler) handle(ctx *gin.Context) {
	alreadyFinished, err := h.commands.OnboardingCompleted(ctx.Request.Context())
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

	domainModel := converter.Wrap(toDomain, requestPayload)
	domainModel.ID = uuid.New()
	domainModel.Enabled = true
	domainModel.Permissions = user.Permissions{
		Hosts:        user.ReadWriteAccessLevel,
		Streams:      user.ReadWriteAccessLevel,
		Certificates: user.ReadWriteAccessLevel,
		Logs:         user.ReadOnlyAccessLevel,
		Integrations: user.ReadWriteAccessLevel,
		AccessLists:  user.ReadWriteAccessLevel,
		Settings:     user.ReadWriteAccessLevel,
		Users:        user.ReadWriteAccessLevel,
		NginxServer:  user.ReadWriteAccessLevel,
		ExportData:   user.ReadOnlyAccessLevel,
		VPNs:         user.ReadWriteAccessLevel,
		Caches:       user.ReadWriteAccessLevel,
	}

	if err = h.commands.Save(ctx.Request.Context(), domainModel, nil); err != nil {
		panic(err)
	}

	usr, err := h.commands.Authenticate(ctx.Request.Context(), domainModel.Username, *domainModel.Password)
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
