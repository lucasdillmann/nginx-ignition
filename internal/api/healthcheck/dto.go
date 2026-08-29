package healthcheck

type statusDTO struct {
	Details []detailDTO `json:"details"`
	Healthy bool        `json:"healthy"`
}

type detailDTO struct {
	Component string `json:"component"`
	Healthy   bool   `json:"healthy"`
}
