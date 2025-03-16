package lifecycle

import "context"

type StartupCommand interface {
	Priority() int
	Async() bool
	Run(ctx context.Context) error
}

type ShutdownCommand interface {
	Priority() int
	Run(ctx context.Context)
}
