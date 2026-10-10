package authorization

type Commands interface {
	JwtSecret() string
}
