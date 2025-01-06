package configuration

type Configuration interface {
	Get(key string) (string, error)
	WithPrefix(prefix string) Configuration
}
