package backup

import "context"

type Commands struct {
	Get func(ctx context.Context) (*Backup, error)
}
