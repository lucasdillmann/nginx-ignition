package host

import "github.com/google/uuid"

type Repository interface {
	ExistsByAccessListId(id uuid.UUID) (bool, error)
}
