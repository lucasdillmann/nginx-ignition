package host

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(id uuid.UUID) (*Host, error)
	DeleteByID(id uuid.UUID) error
	Save(host *Host) error
	FindPage(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[Host], error)
	FindAllEnabled() ([]*Host, error)
	FindDefault() (*Host, error)
	ExistsByID(id uuid.UUID) (bool, error)
	ExistsByCertificateID(certificateId uuid.UUID) (bool, error)
	ExistsByAccessListID(accessListId uuid.UUID) (bool, error)
}
