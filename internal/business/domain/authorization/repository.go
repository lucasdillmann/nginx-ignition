package authorization

import (
	"context"
)

type Repository interface {
	FindJwtSecret(ctx context.Context) (*string, error)
	SaveJwtSecretIfNotExists(ctx context.Context, secret *string) (*string, error)
}
