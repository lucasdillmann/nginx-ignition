package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type onboardingStatusHandler struct {
	command *user.OnboardingCompletedCommand
}

func (h onboardingStatusHandler) handle(context *gin.Context) {
	finished, err := (*h.command)()
	if err != nil {
		panic(err)
	}

	payload := &userOnboardingStatusResponseDto{finished}
	context.JSON(http.StatusOK, payload)
}
