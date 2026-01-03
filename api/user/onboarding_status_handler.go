package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/user"
)

type onboardingStatusHandler struct {
	commands user.Commands
}

func (h onboardingStatusHandler) handle(ctx *gin.Context) {
	finished, err := h.commands.OnboardingCompleted(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	payload := &userOnboardingStatusResponseDTO{finished}
	ctx.JSON(http.StatusOK, payload)
}
