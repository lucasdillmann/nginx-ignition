package main

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func runApplication(lifecycle *lifecycle.Lifecycle) error {
	if err := lifecycle.FireStartup(); err != nil {
		return err
	}

	waitForShutdownSignal(lifecycle)
	log.Println("Shutdown complete")

	return nil
}

func waitForShutdownSignal(lifecycle *lifecycle.Lifecycle) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	<-channel

	log.Println("Application shutdown signal received. Starting graceful shutdown.")
	lifecycle.FireShutdown()
}
