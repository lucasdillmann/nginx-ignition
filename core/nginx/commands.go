package nginx

import (
	"context"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/logline"
)

type LogSearch struct {
	Query            string
	SurroundingLines int
}

type GetConfigFilesInput struct {
	BasePath   string
	ConfigPath string
	LogPath    string
	CachePath  string
	TempPath   string
}

type Commands interface {
	GetHostLogs(
		ctx context.Context,
		hostID uuid.UUID,
		qualifier string,
		lines int,
		search *LogSearch,
	) ([]logline.LogLine, error)
	GetMainLogs(ctx context.Context, lines int, search *LogSearch) ([]logline.LogLine, error)
	GetStatus(ctx context.Context) bool
	GetConfigFiles(ctx context.Context, input GetConfigFilesInput) ([]byte, error)
	GetMetadata(ctx context.Context) (*Metadata, error)
	Reload(ctx context.Context, failIfNotRunning bool) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
