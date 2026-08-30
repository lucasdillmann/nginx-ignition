package server

import (
	"github.com/gin-gonic/gin"

	"github.com/lucasdillmann/nginx-ignition/internal/api/common/apierror"
	"github.com/lucasdillmann/nginx-ignition/internal/api/common/authorization"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/configuration"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/container"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/core/user"
)

func Install() error {
	if err := container.Provide(build); err != nil {
		return err
	}

	return container.Run(registerStartup, registerShutdown)
}

func build(
	cfg *configuration.Configuration,
	userCommands user.Commands,
	i18nCommands i18n.Commands,
) (
	*gin.Engine,
	*state,
	*authorization.ABAC,
	error,
) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(i18nMiddleware(i18nCommands))
	engine.Use(gin.CustomRecoveryWithWriter(nil, apierror.Handler))

	authorizer, err := authorization.New(cfg, userCommands)
	if err != nil {
		return nil, nil, nil, err
	}

	engine.Use(authorizer.HandleRequest)

	return engine, newState(engine), authorizer, nil
}
