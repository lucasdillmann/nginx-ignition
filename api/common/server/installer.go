package server

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/api/common/authorization"
	"dillmann.com.br/nginx-ignition/core/common/configuration"
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func Install(container *dig.Container) error {
	if err := container.Provide(build); err != nil {
		return err
	}

	if err := container.Invoke(registerStartup); err != nil {
		return err
	}

	if err := container.Invoke(registerShutdown); err != nil {
		return err
	}

	return nil
}

func build(
	configuration *configuration.Configuration,
	repository user.Repository,
) (
	*gin.Engine,
	*state,
	*authorization.ABAC,
	error,
) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.CustomRecoveryWithWriter(nil, api_error.Handler))

	authorizer, err := authorization.New(configuration, repository)
	if err != nil {
		return nil, nil, nil, err
	}

	engine.Use(authorizer.HandleRequest)

	return engine, newState(engine), authorizer, nil
}
