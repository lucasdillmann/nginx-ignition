package main

import (
	"dillmann.com.br/nginx-ignition/core"
	"dillmann.com.br/nginx-ignition/database"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	installAllModulesBeans(container)
}

func installAllModulesBeans(container *dig.Container) {
	if err := installBeans(container); err != nil {
		panic(err)
	}

	if err := database.InstallBeans(container); err != nil {
		panic(err)
	}

	if err := core.InstallBeans(container); err != nil {
		panic(err)
	}
}
