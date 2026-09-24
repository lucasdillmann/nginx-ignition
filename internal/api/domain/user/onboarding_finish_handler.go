package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/api/core/converter"

	user2 "github.com/lucasdillmann/nginx-ignition/internal/business/domain/user"
)

type onboardingFinishHandler struct {
	commands   user2.Commands
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
	domainModel.Permissions = user2.Permissions{
		Hosts:        user2.ReadWriteAccessLevel,
		Streams:      user2.ReadWriteAccessLevel,
		Certificates: user2.ReadWriteAccessLevel,
		Logs:         user2.ReadOnlyAccessLevel,
		Integrations: user2.ReadWriteAccessLevel,
		AccessLists:  user2.ReadWriteAccessLevel,
		Settings:     user2.ReadWriteAccessLevel,
		Users:        user2.ReadWriteAccessLevel,
		NginxServer:  user2.ReadWriteAccessLevel,
		ExportData:   user2.ReadOnlyAccessLevel,
		VPNs:         user2.ReadWriteAccessLevel,
		Caches:       user2.ReadWriteAccessLevel,
		TrafficStats: user2.ReadOnlyAccessLevel,
	}

	if err = h.commands.FinishOnboarding(ctx.Request.Context(), domainModel); err != nil {
		if errors.Is(err, user2.ErrOnboardingAlreadyCompleted) {
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

	if outcome != user2.AuthenticationSuccessful || usr == nil {
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
