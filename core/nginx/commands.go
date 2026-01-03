package nginx

import (
	"context"

	"github.com/google/uuid"
)

type GetConfigFilesInput struct {
	BasePath   string
	ConfigPath string
	LogPath    string
	CachePath  string
}

type Commands interface {
	GetHostLogs(
		ctx context.Context,
		hostID uuid.UUID,
		qualifier string,
		lines int,
	) ([]string, error)
	GetMainLogs(ctx context.Context, lines int) ([]string, error)
	GetStatus(ctx context.Context) bool
	GetConfigFiles(ctx context.Context, input GetConfigFilesInput) ([]byte, error)
	GetMetadata(ctx context.Context) (*Metadata, error)
	Reload(ctx context.Context, failIfNotRunning bool) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
