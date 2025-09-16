package backup

import "context"

type Repository interface {
	Get(ctx context.Context) (*Backup, error)
}
