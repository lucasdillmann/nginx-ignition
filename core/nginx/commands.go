package nginx

import (
	"context"

	"github.com/google/uuid"
)

type Commands struct {
	GetHostLogs    func(ctx context.Context, hostId uuid.UUID, qualifier string, lines int) ([]string, error)
	GetMainLogs    func(ctx context.Context, lines int) ([]string, error)
	GetStatus      func(ctx context.Context) bool
	GetConfigFiles func(ctx context.Context) ([]byte, error)
	Reload         func(ctx context.Context, failIfNotRunning bool) error
	Start          func(ctx context.Context) error
	Stop           func(ctx context.Context) error
}
