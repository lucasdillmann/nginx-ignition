package frontend

type configurationDto struct {
	CodeEditor codeEditorDto `json:"codeEditor"`
	Version    versionDto    `json:"version"`
}

type codeEditorDto struct {
	ApiKey *string `json:"apiKey"`
}

type versionDto struct {
	Current *string `json:"current"`
	Latest  *string `json:"latest"`
}
