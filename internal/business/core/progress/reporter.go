package progress

type Reporter interface {
	Report(Event)
}
