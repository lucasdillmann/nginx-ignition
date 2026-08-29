package main

import (
	"nginx-ignition/internal/application"
	"nginx-ignition/internal/core/common/log"
	"nginx-ignition/internal/core/common/version"
)

func main() {
	log.Infof("Welcome to nginx ignition %s", version.Number)

	if err := application.Start(); err != nil {
		log.Fatalf("Application failed to start: %s", err)
	}
}
