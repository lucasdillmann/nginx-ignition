package backup

import (
	"context"
)

type service struct {
	repository Repository
}

func newService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) get(ctx context.Context) (*Backup, error) {
	return s.repository.Get(ctx)
}
