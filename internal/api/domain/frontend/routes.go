package frontend

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/core/authorization"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
)

func Install(
	router *gin.Engine,
	cfg *configuration.Configuration,
	authorizer *authorization.ABAC,
) {
	settingsHandlerInstance := &configurationHandler{cfg}
	router.Handle(http.MethodGet, "/api/frontend/configuration", settingsHandlerInstance.handle)
	authorizer.AllowAllUsers(http.MethodGet, "/api/frontend/configuration")

	router.NoRoute(newStaticFilesHandler().handle)
}

func newStaticFilesHandler() staticFilesHandler {
	files, err := fs.Sub(staticAssets, "static")
	if err != nil {
		log.Warnf("Unable to load embedded frontend assets: %s", err)
		return staticFilesHandler{}
	}

	if _, err := fs.Stat(files, indexFile); err != nil {
		log.Warnf(
			"Frontend assets are not embedded in this binary. Build the frontend and rebuild the application.",
		)
		return staticFilesHandler{}
	}

	return staticFilesHandler{files}
}
