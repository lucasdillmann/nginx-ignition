package lifecycle

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/log"
	"sort"
)

type Lifecycle struct {
	startupCommands  []StartupCommand
	shutdownCommands []ShutdownCommand
}

func New() *Lifecycle {
	return &Lifecycle{
		startupCommands:  []StartupCommand{},
		shutdownCommands: []ShutdownCommand{},
	}
}

func (l *Lifecycle) RegisterStartup(command StartupCommand) {
	l.startupCommands = append(l.startupCommands, command)
}

func (l *Lifecycle) RegisterShutdown(command ShutdownCommand) {
	l.shutdownCommands = append(l.shutdownCommands, command)
}

func (l *Lifecycle) FireStartup(ctx context.Context) error {
	sort.Slice(l.startupCommands, func(left, right int) bool {
		leftCommand := l.startupCommands[left]
		rightCommand := l.startupCommands[right]
		return leftCommand.Priority() < rightCommand.Priority()
	})

	for _, command := range l.startupCommands {
		if command.Async() {
			go func() {
				if err := command.Run(ctx); err != nil {
					log.Warnf("Startup task failed: %s", err)
				}
			}()
		} else {
			if err := command.Run(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Lifecycle) FireShutdown(ctx context.Context) {
	sort.Slice(l.shutdownCommands, func(left, right int) bool {
		leftCommand := l.shutdownCommands[left]
		rightCommand := l.shutdownCommands[right]
		return leftCommand.Priority() < rightCommand.Priority()
	})

	for _, command := range l.shutdownCommands {
		command.Run(ctx)
	}
}
