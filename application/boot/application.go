package boot

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartApplication() error {
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	container, err := startContainer()
	if err != nil {
		return err
	}

	err = container.Invoke(func(lifecycle *lifecycle.Lifecycle) error {
		return startLifecycle(lifecycle, startTime)
	})
	if err != nil {
		return err
	}

	return nil
}

func startLifecycle(lifecycle *lifecycle.Lifecycle, startTime int64) error {
	if err := lifecycle.FireStartup(); err != nil {
		return err
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Application started in %d ms", endTime-startTime)

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
