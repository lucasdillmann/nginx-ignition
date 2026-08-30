package frontend

type configurationDTO struct {
	Version versionDTO `json:"version"`
}

type versionDTO struct {
	Current *string `json:"current"`
	Latest  *string `json:"latest"`
}
