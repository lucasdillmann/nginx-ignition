package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type onboardingStatusHandler struct {
	commands *user.Commands
}

func (h onboardingStatusHandler) handle(ctx *gin.Context) {
	finished, err := h.commands.OnboardingCompleted(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	payload := &userOnboardingStatusResponseDto{finished}
	ctx.JSON(http.StatusOK, payload)
}
