package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/api/common/converter"
	"nginx-ignition/internal/core/user"
)

type onboardingFinishHandler struct {
	commands   user.Commands
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

	requestPayload := &userRequestDTO{}
	if err = ctx.BindJSON(requestPayload); err != nil {
		panic(err)
	}

	domainModel := converter.Wrap(ctx.Request.Context(), toDomain, requestPayload)
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
		TrafficStats: user.ReadOnlyAccessLevel,
	}

	if err = h.commands.FinishOnboarding(ctx.Request.Context(), domainModel); err != nil {
		if errors.Is(err, user.ErrOnboardingAlreadyCompleted) {
			ctx.Status(http.StatusForbidden)
			return
		}

		panic(err)
	}

	outcome, usr, err := h.commands.Authenticate(
		ctx.Request.Context(),
		domainModel.Username,
		*domainModel.Password,
		"",
	)
	if err != nil {
		panic(err)
	}

	if outcome != user.AuthenticationSuccessful || usr == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	token, err := h.authorizer.Jwt().GenerateToken(usr)
	if err != nil {
		panic(err)
	}

	responsePayload := &userLoginResponseDTO{*token}
	ctx.JSON(http.StatusOK, responsePayload)
}
