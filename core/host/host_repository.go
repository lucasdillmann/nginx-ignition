package host

import "github.com/google/uuid"

type HostRepository interface {
	ExistsByAccessListId(id uuid.UUID) (bool, error)
}
