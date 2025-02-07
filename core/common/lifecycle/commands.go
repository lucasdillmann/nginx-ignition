package lifecycle

type StartupCommand interface {
	Priority() int
	Async() bool
	Run() error
}

type ShutdownCommand interface {
	Priority() int
	Run()
}
