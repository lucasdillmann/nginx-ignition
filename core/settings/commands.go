package settings

import "context"

type Commands struct {
	Get  func(ctx context.Context) (*Settings, error)
	Save func(ctx context.Context, settings *Settings) error
}
