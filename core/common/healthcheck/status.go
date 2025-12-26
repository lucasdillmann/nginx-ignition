package healthcheck

type Status struct {
	Healthy bool
	Details []Detail
}

type Detail struct {
	ID    string
	Error error
}
