package main

import (
	"github.com/lucasdillmann/nginx-ignition/internal/application"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/log"
	"github.com/lucasdillmann/nginx-ignition/internal/core/common/version"
)

func main() {
	log.Infof("Welcome to nginx ignition %s", version.Number)

	if err := application.Start(); err != nil {
		log.Fatalf("Application failed to start: %s", err)
	}
}
