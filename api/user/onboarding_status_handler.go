package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type onboardingStatusHandler struct {
	command *user.OnboardingCompletedCommand
}

func (h onboardingStatusHandler) handle(ctx *gin.Context) {
	finished, err := (*h.command)(ctx.Request.Context())
	if err != nil {
		panic(err)
	}

	payload := &userOnboardingStatusResponseDto{finished}
	ctx.JSON(http.StatusOK, payload)
}
