package nginx

import (
	"github.com/google/uuid"
)

type GetHostLogsCommand = func(hostId uuid.UUID, qualifier string, lines int) ([]string, error)

type GetMainLogsCommand = func(lines int) ([]string, error)

type GetStatusCommand = func() bool

type ReloadCommand = func(failIfNotRunning bool) error

type StartCommand = func() error

type StopCommand = func(pid *int) error
