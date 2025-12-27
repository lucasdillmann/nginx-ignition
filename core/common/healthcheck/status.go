package healthcheck

type Status struct {
	Details []Detail
	Healthy bool
}

type Detail struct {
	Error error
	ID    string
}
