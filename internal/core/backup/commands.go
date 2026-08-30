package backup

import "context"

type Commands interface {
	Get(ctx context.Context) (*Backup, error)
}
