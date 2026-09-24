package logline

type ResponseDTO struct {
	Highlight  *HighlightResponseDTO `json:"highlight,omitempty"`
	Contents   string                `json:"contents"`
	LineNumber int                   `json:"lineNumber"`
}

type HighlightResponseDTO struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
