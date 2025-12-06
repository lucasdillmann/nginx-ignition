package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/core/common/healthcheck"
)

type livenessHandler struct {
	healthCheck *healthcheck.HealthCheck
}

func (h livenessHandler) handle(ctx *gin.Context) {
	status := h.healthCheck.Status(ctx.Request.Context())

	payload := toDto(status)
	statusCode := http.StatusOK

	if !status.Healthy {
		statusCode = http.StatusServiceUnavailable
	}

	ctx.JSON(statusCode, payload)
}
