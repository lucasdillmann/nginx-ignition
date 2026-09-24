package logline

type LogLine struct {
	Highlight  *Highlight
	Contents   string
	LineNumber int
}

type Highlight struct {
	Start int
	End   int
}
