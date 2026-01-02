package nginx

import (
	"context"
	"fmt"
	"time"

	"dillmann.com.br/nginx-ignition/core/common/lifecycle"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/core/common/retry"
)

type startup struct {
	service       *service
	retryAttempts int
	retryDelay    time.Duration
}

func registerStartup(lc *lifecycle.Lifecycle, service *service) {
	instance := startup{
		service:       service,
		retryAttempts: 15,
		retryDelay:    time.Second * 5,
	}

	lc.RegisterStartup(instance)
}

func (s startup) Run(ctx context.Context) error {
	go s.service.attachListeners()

	metadata, err := s.service.getMetadata(ctx)
	if err != nil || metadata == nil {
		log.Warnf("Unable to detect nginx binary details: %v", err)
	} else {
		log.Infof(
			"nginx detected with version %s, %s stream support, %s TLS SNI support, and %s code execution support",
			metadata.Version,
			metadata.StreamSupportType(),
			metadata.SNISupportType(),
			metadata.RunCodeSupportType(),
		)
	}

	autoRetry := &retry.Retry{
		Action:               func() error { return s.service.start(ctx) },
		Callback:             s.handleRetryCallback,
		Attempts:             s.retryAttempts,
		DelayBetweenAttempts: s.retryDelay,
	}

	return autoRetry.Start()
}

func (s startup) handleRetryCallback(err error, attempt int, completed bool) {
	if err == nil {
		return
	}

	var msg string
	if completed {
		msg = fmt.Sprintf("Unable to start the nginx server (no new retries will be made): %v", err)
	} else {
		msg = fmt.Sprintf(
			"Unable to start the nginx server (retrying in %.0f seconds; attempt %d of %d): %v",
			s.retryDelay.Seconds(),
			attempt+1,
			s.retryAttempts,
			err,
		)
	}

	log.Warn(msg)
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return true
}
