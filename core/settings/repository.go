package settings

type Repository interface {
	Get() (*Settings, error)
	Save(settings *Settings) error
}
