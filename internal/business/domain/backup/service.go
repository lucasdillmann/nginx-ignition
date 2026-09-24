package backup

import (
	"context"
)

type service struct {
	repository Repository
}

func newCommands(repository Repository) Commands {
	return &service{
		repository: repository,
	}
}

func (s *service) Get(ctx context.Context) (*Backup, error) {
	return s.repository.Get(ctx)
}
