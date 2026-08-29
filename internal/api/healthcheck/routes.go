package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/healthcheck"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
)

const (
	apiPath = "/api/health"
)

func Install(
	router *gin.Engine,
	authorizer *authorization.ABAC,
	healthCheck *healthcheck.HealthCheck,
	cfg *configuration.Configuration,
) {
	enabled, err := cfg.GetBoolean("nginx-ignition.health-check.enabled")
	if err != nil {
		log.Warnf(
			"Unable to check if health check endpoints should be enabled (%v). Keeping them disabled as a fallback.",
			err,
		)
		return
	}

	if !enabled {
		log.Warnf("Health check endpoints disabled by configuration")
		return
	}

	basePath := router.Group(apiPath)
	basePath.GET("/liveness", livenessHandler{healthCheck}.handle)
	basePath.GET("/readiness", readinessHandler{}.handle)

	authorizer.AllowAnonymous(http.MethodGet, apiPath+"/liveness")
	authorizer.AllowAnonymous(http.MethodGet, apiPath+"/readiness")
}
