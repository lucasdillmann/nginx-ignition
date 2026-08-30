package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/healthcheck"
)

type livenessHandler struct {
	healthCheck *healthcheck.HealthCheck
}

func (h livenessHandler) handle(ctx *gin.Context) {
	status := h.healthCheck.Status(ctx.Request.Context())

	payload := toDTO(status)
	statusCode := http.StatusOK

	if !status.Healthy {
		statusCode = http.StatusServiceUnavailable
	}

	ctx.JSON(statusCode, payload)
}
