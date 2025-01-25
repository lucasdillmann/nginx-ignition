package integration

type Repository interface {
	FindByID(id string) (*Integration, error)
	Save(integration *Integration) error
}
