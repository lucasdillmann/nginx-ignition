package settings

type GetCommand func() (*Settings, error)

type SaveCommand func(settings *Settings) error
