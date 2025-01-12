package server

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
	"dillmann.com.br/nginx-ignition/api/common/authentication"
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
) (*gin.Engine, *state, *authentication.Middleware, error) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.CustomRecoveryWithWriter(nil, api_error.Handler))

	authenticationMiddleware, err := authentication.New(configuration, &repository)
	if err != nil {
		return nil, nil, nil, err
	}

	engine.Use(authenticationMiddleware.HandleRequest)

	return engine, &state{engine: engine}, authenticationMiddleware, nil
}
