package integration

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Integration, error)
	Save(ctx context.Context, integration *Integration) error
}
