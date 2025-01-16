package nginx

import "dillmann.com.br/nginx-ignition/core/common/lifecycle"

type startup struct {
	command StartCommand
}

func registerStartup(lifecycle *lifecycle.Lifecycle, command StartCommand) {
	lifecycle.RegisterStartup(startup{command})
}

func (s startup) Run() error {
	return s.command()
}

func (s startup) Priority() int {
	return startupPriority
}

func (s startup) Async() bool {
	return true
}
