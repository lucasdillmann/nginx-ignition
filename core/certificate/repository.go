package certificate

import (
	"dillmann.com.br/nginx-ignition/core/common/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(id uuid.UUID) (*Certificate, error)
	ExistsByID(id uuid.UUID) (bool, error)
	DeleteByID(id uuid.UUID) error
	Save(certificate *Certificate) error
	FindPage(pageSize, pageNumber int, searchTerms *string) (*pagination.Page[*Certificate], error)
	FindAllDueToRenew() ([]*Certificate, error)
}
