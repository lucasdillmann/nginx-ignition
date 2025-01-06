package main

import (
	"dillmann.com.br/nginx-ignition/core"
	"dillmann.com.br/nginx-ignition/database"
	"go.uber.org/dig"
)

func main() {
	digContainer := dig.New()
	installApplicationModules(digContainer)
}

func installApplicationModules(container *dig.Container) {
	if err := database.RegisterDatabaseBeans(container); err != nil {
		panic(err)
	}

	if err := core.RegisterCoreBeans(container); err != nil {
		panic(err)
	}
}
