package server

import (
	"github.com/gin-gonic/gin"

	"nginx-ignition/internal/api/common/apierror"
	"nginx-ignition/internal/api/common/authorization"
	"nginx-ignition/internal/core/common/configuration"
	"nginx-ignition/internal/core/common/container"
	"nginx-ignition/internal/core/common/i18n"
	"nginx-ignition/internal/core/user"
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
