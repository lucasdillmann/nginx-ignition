package main

import (
	"dillmann.com.br/nginx-ignition/application/boot"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/version"
)

func main() {
	if err := log.Init(); err != nil {
		panic(err)
	}

	log.Infof("Welcome to nginx ignition %s", version.Number)

	if err := boot.StartApplication(); err != nil {
		log.Fatalf("Application failed to start: %s", err)
	}
}
