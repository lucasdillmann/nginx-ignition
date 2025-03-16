package nginx

import (
	"context"
	"github.com/google/uuid"
)

type GetHostLogsCommand = func(ctx context.Context, hostId uuid.UUID, qualifier string, lines int) ([]string, error)

type GetMainLogsCommand = func(ctx context.Context, lines int) ([]string, error)

type GetStatusCommand = func(ctx context.Context) bool

type ReloadCommand = func(ctx context.Context, failIfNotRunning bool) error

type StartCommand = func(ctx context.Context) error

type StopCommand = func(ctx context.Context, pid *int) error
