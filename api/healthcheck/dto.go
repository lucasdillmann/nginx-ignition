package healthcheck

type statusDto struct {
	Healthy bool        `json:"healthy"`
	Details []detailDto `json:"details"`
}

type detailDto struct {
	Component string `json:"component"`
	Healthy   bool   `json:"healthy"`
}
