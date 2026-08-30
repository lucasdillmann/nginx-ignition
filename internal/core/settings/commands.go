package settings

import "context"

type Commands interface {
	Get(ctx context.Context) (*Settings, error)
	Save(ctx context.Context, settings *Settings) error
}
