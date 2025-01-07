package main

import (
	"dillmann.com.br/nginx-ignition/core"
	"dillmann.com.br/nginx-ignition/database"
	"go.uber.org/dig"
	"log"
)

func main() {
	log.Println("Welcome to nginx ignition")

	container := dig.New()
	installAllModulesBeans(container)

	log.Fatal("Application lifecycle isn't implemented yet")
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
