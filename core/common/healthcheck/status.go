package healthcheck

type Status struct {
	Healthy bool
	Details []*Detail
}

type Detail struct {
	ID      string
	Healthy bool
	Error   error
}
