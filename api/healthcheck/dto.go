package healthcheck

type statusDto struct {
	Details []detailDto `json:"details"`
	Healthy bool        `json:"healthy"`
}

type detailDto struct {
	Component string `json:"component"`
	Healthy   bool   `json:"healthy"`
}
