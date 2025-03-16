package frontend

type configurationDto struct {
	CodeEditor codeEditorDto `json:"codeEditor"`
}

type codeEditorDto struct {
	ApiKey *string `json:"apiKey"`
}
