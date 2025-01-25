package integration

type Repository interface {
	FindById(id string) (*Integration, error)
	Save(integration *Integration) error
}
