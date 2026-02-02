package logline

import (
	"dillmann.com.br/nginx-ignition/core/nginx"
)

func ToResponseDTOs(logLines []nginx.LogLine) []ResponseDTO {
	result := make([]ResponseDTO, len(logLines))
	for index, logLine := range logLines {
		result[index] = ToResponseDTO(logLine)
	}

	return result
}

func ToResponseDTO(logLine nginx.LogLine) ResponseDTO {
	var highlight *HighlightResponseDTO
	if logLine.Highlight != nil {
		highlight = &HighlightResponseDTO{
			Start: logLine.Highlight.Start,
			End:   logLine.Highlight.End,
		}
	}

	return ResponseDTO{
		LineNumber: logLine.LineNumber,
		Contents:   logLine.Contents,
		Highlight:  highlight,
	}
}
