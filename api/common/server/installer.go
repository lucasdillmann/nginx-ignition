package server

import (
	"github.com/gin-gonic/gin"

	"dillmann.com.br/nginx-ignition/api/common/apierror"
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/user"
)

func Install() error {
	if err := container.Provide(build); err != nil {
		return err
	}

	return container.Run(registerStartup, registerShutdown)
}

func build(
	cfg *configuration.Configuration,
	commands *user.Commands,
) (
	*gin.Engine,
	*state,
	*authorization.ABAC,
	error,
) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.CustomRecoveryWithWriter(nil, apierror.Handler))

	authorizer, err := authorization.New(cfg, commands)
	if err != nil {
		return nil, nil, nil, err
	}

	engine.Use(authorizer.HandleRequest)

	return engine, newState(engine), authorizer, nil
}
