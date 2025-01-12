package server

import (
	"dillmann.com.br/nginx-ignition/api/common/api_error"
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

func build() (*gin.Engine, *state) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(
		gin.CustomRecoveryWithWriter(nil, api_error.Handler),
	)

	return engine, &state{engine: engine}
}
