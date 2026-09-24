package main

import (
	"github.com/lucasdillmann/nginx-ignition/internal/application"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/log"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/version"
)

func main() {
	log.Infof("Welcome to nginx ignition %s", version.Number)

	if err := application.Start(); err != nil {
		log.Fatalf("Application failed to start: %s", err)
	}
}
