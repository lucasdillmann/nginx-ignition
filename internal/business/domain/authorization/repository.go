package authorization

import (
	"context"
)

type Repository interface {
	FindJwtSecret(ctx context.Context) (*string, error)
	SaveJwtSecret(ctx context.Context, secret *string) error
}
