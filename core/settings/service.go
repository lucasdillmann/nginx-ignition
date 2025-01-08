package settings

type service struct {
	repository *Repository
}

func (s *service) Get() (*Settings, error) {
	return (*s.repository).Get()
}

func (s *service) Save(settings *Settings) error {
	return (*s.repository).Save(settings)
}
