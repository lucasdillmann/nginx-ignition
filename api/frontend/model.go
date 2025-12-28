package frontend

type configurationDTO struct {
	CodeEditor codeEditorDTO `json:"codeEditor"`
	Version    versionDTO    `json:"version"`
}

type codeEditorDTO struct {
	APIKey *string `json:"apiKey"`
}

type versionDTO struct {
	Current *string `json:"current"`
	Latest  *string `json:"latest"`
}
