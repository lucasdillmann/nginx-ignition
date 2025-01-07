package lifecycle

import (
	"sort"
)

type Lifecycle struct {
	startupCommands  *[]StartupCommand
	shutdownCommands *[]ShutdownCommand
}

func New() *Lifecycle {
	return &Lifecycle{
		startupCommands:  &[]StartupCommand{},
		shutdownCommands: &[]ShutdownCommand{},
	}
}

func (l *Lifecycle) RegisterStartup(command StartupCommand) {
	updatedValues := append(*l.startupCommands, command)
	l.startupCommands = &updatedValues
}

func (l *Lifecycle) RegisterShutdown(command *ShutdownCommand) {
	updatedValues := append(*l.shutdownCommands, *command)
	l.shutdownCommands = &updatedValues
}

func (l *Lifecycle) FireStartup() error {
	commands := *l.startupCommands
	sort.Slice(commands, func(left, right int) bool {
		return commands[left].Priority() < commands[right].Priority()
	})

	for _, command := range commands {
		if command.Async() {
			go func() { _ = command.Run() }()
		} else {
			if err := command.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Lifecycle) FireShutdown() {
	commands := *l.shutdownCommands
	sort.Slice(commands, func(left, right int) bool {
		return commands[left].Priority() < commands[right].Priority()
	})

	for _, command := range commands {
		command.Run()
	}
}
