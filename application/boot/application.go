package boot

import (
	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
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

	return container.Invoke(func(lifecycle *lifecycle.Lifecycle) error {
		return runLifecycle(lifecycle, startTime)
	})
}

func runLifecycle(lifecycle *lifecycle.Lifecycle, startTime int64) error {
	if err := lifecycle.FireStartup(); err != nil {
		return err
	}

	endTime := time.Now().UnixNano() / int64(time.Millisecond)
	log.Infof("Application started in %d ms", endTime-startTime)

	receivedSignal := waitForShutdownSignal()

	log.Infof("Application shutdown signal received (%s). Starting graceful shutdown.", receivedSignal)
	lifecycle.FireShutdown()

	log.Infof("Shutdown complete")
	return nil
}

func waitForShutdownSignal() os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV)

	return <-channel
}
