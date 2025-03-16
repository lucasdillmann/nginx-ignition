package settings

import "context"

type GetCommand func(ctx context.Context) (*Settings, error)

type SaveCommand func(ctx context.Context, settings *Settings) error
